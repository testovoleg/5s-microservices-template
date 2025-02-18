package kafka

import (
	"context"
	"sync"

	"github.com/go-playground/validator"
	"github.com/segmentio/kafka-go"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/service"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
)

const (
	PoolSize = 30
)

type coreMessageProcessor struct {
	log     logger.Logger
	cfg     *config.Config
	v       *validator.Validate
	ps      *service.Service
	metrics *metrics.CoreServiceMetrics
}

func NewCoreMessageProcessor(log logger.Logger, cfg *config.Config, v *validator.Validate, ps *service.Service, metrics *metrics.CoreServiceMetrics) *coreMessageProcessor {
	return &coreMessageProcessor{log: log, cfg: cfg, v: v, ps: ps, metrics: metrics}
}

func (s *coreMessageProcessor) ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := r.FetchMessage(ctx)
		if err != nil {
			s.log.Warnf("workerID: %v, err: %v", workerID, err)
			continue
		}

		s.logProcessMessage(m, workerID)

		switch m.Topic {
		case s.cfg.KafkaTopics.ProductCreated.TopicName:
			s.processProductCreated(ctx, r, m)
		case s.cfg.KafkaTopics.ProductUpdated.TopicName:
			s.processProductUpdated(ctx, r, m)
		case s.cfg.KafkaTopics.ProductDeleted.TopicName:
			s.processProductDeleted(ctx, r, m)
		}
	}
}
