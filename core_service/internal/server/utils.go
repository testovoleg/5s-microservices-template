package server

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/heptiolabs/healthcheck"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/testovoleg/5s-microservice-template/pkg/constants"
	kafkaClient "github.com/testovoleg/5s-microservice-template/pkg/kafka"
)

const (
	stackSize = 1 << 10 // 1 KB
	// probeTimeout — предел одной readiness-проверки Kafka (и Postgres, если подключён).
	probeTimeout = 5 * time.Second
)

func (s *server) connectKafkaBrokers(ctx context.Context) error {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, s.cfg.Kafka)
	if err != nil {
		return errors.Wrap(err, "kafka.NewKafkaCon")
	}

	s.kafkaConn = kafkaConn

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return errors.Wrap(err, "kafkaConn.Brokers")
	}

	s.log.Infof("kafka connected to brokers: %+v", brokers)

	return nil
}

func (s *server) initKafkaTopics(ctx context.Context) {
	controller, err := s.kafkaConn.Controller()
	if err != nil {
		s.log.WarnMsg("kafkaConn.Controller", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	s.log.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		s.log.WarnMsg("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errcheck

	s.log.Infof("established new kafka controller connection: %s", controllerURI)

	// Новые топики добавлять в оба списка: здесь (создание) и в getConsumerGroupTopics
	// (подписка). Пустое имя (топик выключен конфигом) пропускается.
	toCreate := make([]kafka.TopicConfig, 0, 1)
	for _, t := range []kafka.TopicConfig{
		kafkaTopic(s.cfg.KafkaTopics.WebhookExample),
	} {
		if t.Topic != "" {
			toCreate = append(toCreate, t)
		}
	}
	if err := conn.CreateTopics(toCreate...); err != nil {
		s.log.WarnMsg("kafkaConn.CreateTopics", err)
		return
	}

	s.log.Infof("kafka topics created or already exists: %+v", toCreate)
}

// kafkaTopic — TopicConfig из конфига с дефолтами: partitions/replicationFactor ≤ 0
// брокер отвергает (InvalidPartitions), а забыть их в yaml нового топика легко.
func kafkaTopic(cfg kafkaClient.TopicConfig) kafka.TopicConfig {
	parts := cfg.Partitions
	if parts <= 0 {
		parts = 1
	}
	rf := cfg.ReplicationFactor
	if rf <= 0 {
		rf = 1
	}
	return kafka.TopicConfig{
		Topic:             cfg.TopicName,
		NumPartitions:     parts,
		ReplicationFactor: rf,
	}
}

func (s *server) getConsumerGroupTopics() []string {
	topics := []string{
		s.cfg.KafkaTopics.WebhookExample.TopicName,
	}
	out := make([]string, 0, len(topics))
	for _, t := range topics {
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func (s *server) runHealthCheck(ctx context.Context) {
	health := healthcheck.NewHandler()

	health.AddLivenessCheck(s.cfg.ServiceName, healthcheck.AsyncWithContext(ctx, func() error {
		return nil
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	health.AddReadinessCheck(constants.Redis, healthcheck.AsyncWithContext(ctx, func() error {
		return s.redisClient.Ping(ctx).Err()
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	// Свежее соединение на каждую проверку: долгоживущий s.kafkaConn после переезда
	// брокера остаётся мёртвым, и под навсегда выпадал бы из Service при живом liveness.
	health.AddReadinessCheck(constants.Kafka, healthcheck.AsyncWithContext(ctx, func() error {
		dialCtx, cancel := context.WithTimeout(ctx, probeTimeout)
		defer cancel()
		conn, err := kafkaClient.NewKafkaConn(dialCtx, s.cfg.Kafka)
		if err != nil {
			return err
		}
		defer conn.Close() // nolint: errcheck
		_ = conn.SetDeadline(time.Now().Add(probeTimeout))
		_, err = conn.Brokers()
		return err
	}, time.Duration(s.cfg.Probes.CheckIntervalSeconds)*time.Second))

	go func() {
		s.log.Infof("Core microservice Kubernetes probes listening on port: %s", s.cfg.Probes.Port)
		if err := http.ListenAndServe(s.cfg.Probes.Port, health); err != nil {
			s.log.WarnMsg("ListenAndServe", err)
		}
	}()
}
