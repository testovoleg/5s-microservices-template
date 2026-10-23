package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	coreService "github.com/testovoleg/5s-microservice-template/core_service/proto"
	"github.com/testovoleg/5s-microservice-template/graphql_service/config"
	graph_resolvers "github.com/testovoleg/5s-microservice-template/graphql_service/internal/app/delivery/graphql"
	"github.com/testovoleg/5s-microservice-template/graphql_service/internal/app/service"
	"github.com/testovoleg/5s-microservice-template/graphql_service/internal/client"
	"github.com/testovoleg/5s-microservice-template/graphql_service/internal/middlewares"
	graph "github.com/testovoleg/5s-microservice-template/graphql_service/schema"
	"github.com/testovoleg/5s-microservice-template/pkg/interceptors"
	"github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"go.opentelemetry.io/otel"
)

const (
	shutdownTimeout = 10 * time.Second
	// graphqlComplexityLimit rejects queries whose calculated field complexity exceeds this value
	// (DoS protection; raise if legitimate nested/list queries start failing).
	graphqlComplexityLimit = 1500
)

type server struct {
	log  logger.Logger
	cfg  *config.Config
	v    *validator.Validate
	mw   middlewares.MiddlewareManager
	im   interceptors.InterceptorManager
	echo *echo.Echo
	bs   *service.GraphQLService
	m    *metrics.MetricsManager
	gql  *handler.Server
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

	var otlShutdown func(context.Context) error
	if s.cfg.OTL.Enable {
		provider, shutdown, err := tracing.NewOTLTracer(ctx, s.cfg.OTL)
		if err != nil {
			return err
		}
		otlShutdown = shutdown
		otel.SetTracerProvider(provider)
	}

	coreServiceConn, err := client.NewCoreServiceConn(ctx, s.cfg, s.im)
	if err != nil {
		return err
	}
	defer coreServiceConn.Close() // nolint: errcheck
	coreClient := coreService.NewCoreServiceClient(coreServiceConn)

	kafkaProducer := kafka.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close() // nolint: errcheck

	s.bs = service.NewGraphQLService(s.log, s.cfg, kafkaProducer, coreClient)

	resolverHandlers := graph_resolvers.NewResolverHandlers(s.log, s.mw, s.cfg, s.bs, s.v, s.m)

	s.gql = handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolverHandlers}))
	s.gql.Use(extension.FixedComplexityLimit(graphqlComplexityLimit))

	go func() {
		if err := s.runGraphQLServer(); err != nil {
			s.log.Errorf(" s.runGraphQLServer: %v", err)
			cancel()
		}
	}()
	s.log.Infof("GraphQL Server is listening on PORT: %s", s.cfg.Http.Port)

	s.runHealthCheck(ctx)
	s.m.NewServer(cancel, stackSize)

	<-ctx.Done()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if err := s.echo.Server.Shutdown(shutdownCtx); err != nil {
		s.log.WarnMsg("echo.Server.Shutdown", err)
	}
	if err := s.m.Shutdown(shutdownCtx); err != nil {
		s.log.WarnMsg("metrics.Shutdown", err)
	}
	if otlShutdown != nil {
		if err := otlShutdown(shutdownCtx); err != nil {
			s.log.WarnMsg("otl.Shutdown", err)
		}
	}

	return nil
}
