package kafka

import (
	"context"
	"time"

	"github.com/avast/retry-go"
	"github.com/segmentio/kafka-go"
	"github.com/testovoleg/5s-microservice-template/core_service/internal/mappers"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
	kafkaMessages "github.com/testovoleg/5s-microservice-template/proto/kafka"
	"google.golang.org/protobuf/proto"
)

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var (
	retryOptions = []retry.Option{retry.Attempts(retryAttempts), retry.Delay(retryDelay), retry.DelayType(retry.BackOffDelay)}
)

func (s *coreMessageProcessor) processWebhookExample(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.Get("WebhookExample", metrics.KAFKA).Inc()

	ctx, span := tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, "coreMessageProcessor.processWebhookExample")
	defer span.End()

	msg := &kafkaMessages.Payload{}
	if err := proto.Unmarshal(m.Value, msg); err != nil {
		s.log.WarnMsg("proto.Unmarshal", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	command := mappers.NewWebhookCommand(msg)
	if err := s.v.StructCtx(ctx, command); err != nil {
		s.log.WarnMsg("validate", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	if err := retry.Do(func() error {
		return s.ps.Commands.Webhook.Handle(ctx, command)
	}, append(retryOptions, retry.Context(ctx))...); err != nil {
		s.log.WarnMsg("Webhook.Handle", err)
		s.commitErrMessage(ctx, r, m)
		return
	}

	s.commitMessage(ctx, r, m)
}
