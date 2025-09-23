package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	readerKafka "github.com/testovoleg/5s-microservice-template/core_service/internal/app/delivery/kafka"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/service"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/interceptors"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	redisClient "github.com/testovoleg/5s-microservice-template/pkg/redis"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"go.opentelemetry.io/otel"
)

type server struct {
	log            logger.Logger
	cfg            *config.Config
	v              *validator.Validate
	kafkaConn      *kafka.Conn
	im             interceptors.InterceptorManager
	redisClient    redis.UniversalClient
	svc            *service.Service
	metrics        *metrics.CoreServiceMetrics
	keycloakClient *gocloak.GoCloak
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{log: log, cfg: cfg, v: validator.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s.im = interceptors.NewInterceptorManager(s.log)
	s.metrics = metrics.NewCoreServiceMetrics(s.cfg)

	s.redisClient = redisClient.NewUniversalRedisClient(s.cfg.Redis)
	defer s.redisClient.Close() // nolint: errcheck
	s.log.Infof("Redis connected: %+v", s.redisClient.PoolStats())

	s.keycloakClient = gocloak.NewClient(s.cfg.Keycloak.Host)

	redisRepo := repository.NewRedisRepository(s.log, s.cfg, s.redisClient)
	cloakRepo := repository.NewIDMRepository(s.log, s.cfg, s.keycloakClient)

	s.svc = service.NewAppService(s.log, s.cfg, redisRepo, cloakRepo)

	s.log.Info("Starting Reader Kafka consumers")
	coreMessageProcessor := readerKafka.NewCoreMessageProcessor(s.log, s.cfg, s.v, s.svc, s.metrics)
	cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.log)
	go cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), readerKafka.PoolSize, coreMessageProcessor.ProcessMessages)
	if err := s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer s.kafkaConn.Close() // nolint: errcheck

	s.runHealthCheck(ctx)
	s.runMetrics(cancel)

	if s.cfg.OTL.Enable {
		provider, shutdown, err := tracing.NewOTLTracer(ctx, s.cfg.OTL)
		if err != nil {
			return err
		}
		defer func() { _ = shutdown(ctx) }()

		otel.SetTracerProvider(provider)
	}

	closeGrpcServer, grpcServer, err := s.newReaderGrpcServer()
	if err != nil {
		return errors.Wrap(err, "NewScmGrpcServer")
	}
	defer closeGrpcServer() // nolint: errcheck

	<-ctx.Done()
	grpcServer.GracefulStop()
	return nil
}
