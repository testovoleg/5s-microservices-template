package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	PublishMessageInTopic(ctx context.Context, topicName string, bytes []byte) error
	// PublishKeyedMessage publishes with a partition key: keyAwareBalancer hashes the key so
	// all messages of one entity land on one partition and are consumed in order.
	PublishKeyedMessage(ctx context.Context, topicName string, key, value []byte) error
	Close() error
}

type producer struct {
	log     logger.Logger
	brokers []string
	w       *kafka.Writer
}

// NewProducer create new kafka producer
func NewProducer(log logger.Logger, brokers []string) *producer {
	return &producer{log: log, brokers: brokers, w: NewWriter(brokers, kafka.LoggerFunc(log.Errorf))}
}

func (p *producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	return p.w.WriteMessages(ctx, msgs...)
}

func (p *producer) PublishMessageInTopic(ctx context.Context, name string, value []byte) error {
	return p.PublishKeyedMessage(ctx, name, nil, value)
}

func (p *producer) PublishKeyedMessage(ctx context.Context, name string, key, value []byte) error {
	if len(key) == 0 && len(value) == 0 {
		return nil
	}

	return p.PublishMessage(ctx, kafka.Message{
		Topic:   name,
		Key:     key,
		Value:   value,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(ctx),
	})
}

func (p *producer) Close() error {
	return p.w.Close()
}
