package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/go-playground/validator"
	"github.com/segmentio/kafka-go"
	"github.com/testovoleg/5s-microservice-template/core_service/config"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/app/service"
	"github.com/testovoleg/5s-microservice-template/pkg/logger"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	"github.com/testovoleg/5s-microservice-template/pkg/utils"
)

const (
	PoolSize = 5
)

type coreMessageProcessor struct {
	log     logger.Logger
	cfg     *config.Config
	v       *validator.Validate
	ps      *service.CoreService
	metrics *metrics.MetricsManager
}

func NewCoreMessageProcessor(log logger.Logger, cfg *config.Config, v *validator.Validate, ps *service.CoreService, metrics *metrics.MetricsManager) *coreMessageProcessor {
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

		msgCtx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "coreMessageProcessor.ProcessMessages")
		// The producer's span/trace may have already finished by the time we get here (fire-and-forget),
		// so record the produce->consume gap explicitly on the span - this is what makes the causal
		// order readable in Tempo even though this span shares the producer's (possibly long-closed) trace.
		utils.Attr(span,
			"kafka.produced_at", m.Time.UTC().Format(time.RFC3339Nano),
			"kafka.consume_delay", time.Since(m.Time).String(),
		)
		s.logProcessMessage(msgCtx, m, workerID)

		switch m.Topic {
		case s.cfg.KafkaTopics.WebhookExample.TopicName:
			s.processWebhookExample(msgCtx, r, m)
		}
		span.End()
	}
}
