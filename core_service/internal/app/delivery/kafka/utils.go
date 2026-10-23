package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/testovoleg/5s-microservice-template/pkg/metrics"
	"github.com/testovoleg/5s-microservice-template/pkg/tracing"
)

func (s *coreMessageProcessor) commitMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.Get("Success", metrics.KAFKA).Inc()

	s.log.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset, tracing.TraceIDFromContext(ctx))

	if err := r.CommitMessages(ctx, m); err != nil {
		s.log.WarnMsg("commitMessage", err)
	}
}

func (s *coreMessageProcessor) logProcessMessage(ctx context.Context, m kafka.Message, workerID int) {
	s.log.KafkaProcessMessage(m.Topic, m.Partition, string(m.Value), workerID, m.Offset, m.Time, tracing.TraceIDFromContext(ctx))
}

func (s *coreMessageProcessor) commitErrMessage(ctx context.Context, r *kafka.Reader, m kafka.Message) {
	s.metrics.Get("Error", metrics.KAFKA).Inc()

	s.log.KafkaLogCommittedMessage(m.Topic, m.Partition, m.Offset, tracing.TraceIDFromContext(ctx))
	if err := r.CommitMessages(ctx, m); err != nil {
		s.log.WarnMsg("commitMessage", err)
	}
}
