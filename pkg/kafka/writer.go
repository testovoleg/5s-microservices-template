package kafka

import (
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

// keyAwareBalancer hashes Kafka keys so messages of one entity (e.g. one company or one
// prefix) stay on a single partition and are consumed in order. Messages without a key
// keep LeastBytes, which spreads bulk work evenly across partitions.
type keyAwareBalancer struct {
	hash     kafka.Hash
	fallback kafka.LeastBytes
}

func (b *keyAwareBalancer) Balance(msg kafka.Message, partitions ...int) int {
	if len(msg.Key) > 0 {
		return b.hash.Balance(msg, partitions...)
	}
	return b.fallback.Balance(msg, partitions...)
}

// NewWriter create new configured kafka writer
func NewWriter(brokers []string, errLogger kafka.Logger) *kafka.Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &keyAwareBalancer{},
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		ErrorLogger:  errLogger,
		Compression:  compress.Snappy,
		ReadTimeout:  writerReadTimeout,
		WriteTimeout: writerWriteTimeout,
		Async:        false,
	}
	return w
}
