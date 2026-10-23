package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	readerKafka "github.com/testovoleg/5s-microservice-template/core_service/internal/app/delivery/kafka"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/repository"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/service"
	"github.com/testovoleg/5s-microservice-template/pkg/interceptors"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	redisClient "github.com/testovoleg/5s-microservice-template/pkg/redis"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"go.opentelemetry.io/otel"
)

// consumerDrainTimeout — сколько ждать остановки Kafka-воркеров при выключении;
// должно быть меньше terminationGracePeriodSeconds пода.
const consumerDrainTimeout = 30 * time.Second

type server struct {
	log            logger.Logger
	cfg            *config.Config
	v              *validator.Validate
	kafkaConn      *kafka.Conn
	im             interceptors.InterceptorManager
	redisClient    redis.UniversalClient
	svc            *service.CoreService
	m              *metrics.MetricsManager
	keycloakClient *gocloak.GoCloak
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{log: log, cfg: cfg, v: validator.New()}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s.im = interceptors.NewInterceptorManager(s.log)
	s.m = metrics.NewMetricsManager(s.log, s.cfg.ServiceName, s.cfg.Probes.PrometheusPath, s.cfg.Probes.PrometheusPort)

	s.redisClient = redisClient.NewUniversalRedisClient(s.cfg.Redis)
	defer s.redisClient.Close() // nolint: errcheck
	s.log.Infof("Redis connected: %+v", s.redisClient.PoolStats())

	kafkaProducer := kafkaClient.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close()

	s.keycloakClient = gocloak.NewClient(s.cfg.Keycloak.Host)

	idmRepo := repository.NewIDMRepository(s.log, s.cfg, s.keycloakClient)
	redisRepo := repository.NewRedisRepository(s.log, s.cfg, s.redisClient)
	adminRepo := repository.NewAdminRepository(s.log, s.cfg)
	storageRepo := repository.NewStorageRepository(s.log, s.cfg)

	s.svc = service.NewCoreService(s.log, s.cfg, idmRepo, adminRepo, redisRepo, kafkaProducer, storageRepo)

	var otlShutdown func(context.Context) error
	if s.cfg.OTL.Enable {
		provider, shutdown, err := tracing.NewOTLTracer(ctx, s.cfg.OTL)
		if err != nil {
			return err
		}
		otlShutdown = shutdown
		otel.SetTracerProvider(provider)
	}

	if err := s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer s.kafkaConn.Close() // nolint: errcheck

	if s.cfg.Kafka.InitTopics {
		s.initKafkaTopics(ctx)
	}

	s.log.Info("Starting Reader Kafka consumers")
	coreMessageProcessor := readerKafka.NewCoreMessageProcessor(s.log, s.cfg, s.v, s.svc, s.m)
	cg := kafkaClient.NewConsumerGroup(s.cfg.Kafka.Brokers, s.cfg.Kafka.GroupID, s.log)
	consumerDone := make(chan struct{})
	go func() {
		defer close(consumerDone)
		cg.ConsumeTopic(ctx, s.getConsumerGroupTopics(), readerKafka.PoolSize, coreMessageProcessor.ProcessMessages)
	}()

	s.runHealthCheck(ctx)
	s.m.NewServer(cancel, stackSize)

	closeGrpcServer, grpcServer, err := s.newCoreGrpcServer()
	if err != nil {
		return errors.Wrap(err, "NewScmGrpcServer")
	}
	defer closeGrpcServer() // nolint: errcheck

	<-ctx.Done()
	grpcServer.GracefulStop()
	// Дождаться воркеров Kafka до закрытия пулов: ConsumeTopic закрывает Reader
	// (LeaveGroup → быстрый ребаланс на другие поды), а незакоммиченные сообщения
	// переиграются там. Иначе defer'ы выше рвали бы Redis/Postgres под живыми воркерами.
	select {
	case <-consumerDone:
	case <-time.After(consumerDrainTimeout):
		s.log.Warnf("kafka consumer did not stop within %s, closing anyway", consumerDrainTimeout)
	}
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
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
