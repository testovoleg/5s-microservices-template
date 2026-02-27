package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/config"
	v1 "github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/app/delivery/http/v1"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/app/service"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/client"
	"github.com/testovoleg/5s-microservice-template/api_gateway_service/internal/middlewares"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/docs"
	"github.com/testovoleg/5s-microservice-template/pkg/interceptors"
	"github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"go.opentelemetry.io/otel"
)

type server struct {
	log  logger.Logger
	cfg  *config.Config
	v    *validator.Validate
	mw   middlewares.MiddlewareManager
	im   interceptors.InterceptorManager
	echo *echo.Echo
	svc  *service.GatewayService
	m    *metrics.MetricsManager
	rdcl *docs.Redocly
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{log: log, cfg: cfg, echo: echo.New(), v: validator.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s.mw = middlewares.NewMiddlewareManager(s.log, s.cfg)
	s.im = interceptors.NewInterceptorManager(s.log)
	s.m = metrics.NewMetricsManager(s.log, s.cfg.ServiceName, s.cfg.Probes.PrometheusPath, s.cfg.Probes.PrometheusPort)

	s.rdcl = docs.NewRedoclyHtml(s.cfg.Http.Title, s.cfg.Http.BasePath)

	coreServiceConn, err := client.NewCoreServiceConn(ctx, s.cfg, s.im)
	if err != nil {
		return err
	}
	defer coreServiceConn.Close() // nolint: errcheck
	coreClient := coreService.NewCoreServiceClient(coreServiceConn)

	kafkaProducer := kafka.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close() // nolint: errcheck

	s.svc = service.NewGatewayService(s.log, s.cfg, kafkaProducer, coreClient)

	appHandlers := v1.NewAppHandlers(s.echo.Group(s.cfg.Http.V1Path), s.log, s.mw, s.cfg, s.svc, s.v, s.m, s.rdcl)
	appHandlers.MapRoutes()

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.log.Errorf(" s.runHttpServer: %v", err)
			cancel()
		}
	}()
	s.log.Infof("API Gateway is listening on PORT: %s", s.cfg.Http.Port)

	s.runHealthCheck(ctx)
	s.m.NewServer(cancel, stackSize)

	if s.cfg.OTL.Enable {
		provider, shutdown, err := tracing.NewOTLTracer(ctx, s.cfg.OTL)
		if err != nil {
			return err
		}
		defer func() { _ = shutdown(ctx) }()

		otel.SetTracerProvider(provider)
	}

	<-ctx.Done()
	if err := s.echo.Server.Shutdown(ctx); err != nil {
		s.log.WarnMsg("echo.Server.Shutdown", err)
	}

	return nil
}
