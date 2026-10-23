package tracing

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/metadata"
)

// bulkTracingKey marks a context so StartSpan stops creating real spans under it.
type bulkTracingKey struct{}

// noopTracer never talks to the real TracerProvider/exporter — Start returns a Span
// whose End/SetStatus/RecordError/SetAttributes are all safe no-ops.
var noopTracer = noop.NewTracerProvider().Tracer("")

// WithBulkTracing marks ctx so every StartSpan call made under it (directly or through
// nested calls, e.g. repository methods) becomes a no-op instead of creating a real span.
//
// Why this exists: every repository/redis/http call is auto-instrumented via StartSpan.
// That's fine for a normal request with a handful of calls, but a tight loop processing
// a large batch turns into tens of thousands of spans in a single trace. At that volume
// the BatchSpanProcessor's queue (pkg/tracing/opentelemetry.go) overflows and silently
// drops spans — and because a loop's outer span always ends (and is enqueued for export)
// after all of its children, it's disproportionately likely to be the one dropped,
// leaving thousands of orphaned children with no parent in Tempo.
//
// Callers doing bulk/batch work should wrap ctx with this once they know the batch is
// large, and rely on span attributes (see utils.Attr) on their own outer span for
// aggregate visibility instead of a span per low-level call.
func WithBulkTracing(ctx context.Context) context.Context {
	return context.WithValue(ctx, bulkTracingKey{}, true)
}

func isBulkTracing(ctx context.Context) bool {
	v, _ := ctx.Value(bulkTracingKey{}).(bool)
	return v
}

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	if !trace.SpanContextFromContext(ctx).IsValid() {
		return ctx, trace.SpanFromContext(ctx)
	}
	if isBulkTracing(ctx) {
		return noopTracer.Start(ctx, name)
	}
	return otel.GetTracerProvider().Tracer("").Start(ctx, name)
}

// TraceIDFromContext returns the current span's trace id, or "" if none.
// Used in access logs so Grafana Loki derived fields can link to Tempo.
func TraceIDFromContext(ctx context.Context) string {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.HasTraceID() {
		return ""
	}
	return sc.TraceID().String()
}

func StartHttpServerTracerSpan(c echo.Context, operationName string) (context.Context, trace.Span) {
	ctx := c.Request().Context()

	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request().Header))

	ctx, serverSpan := otel.GetTracerProvider().Tracer("").Start(ctx, operationName)

	traceID := serverSpan.SpanContext().TraceID()
	if traceID.IsValid() {
		c.Set("traceid", traceID.String())
		c.Response().Header().Set("x-trace-id", traceID.String())
	}

	return ctx, serverSpan
}

func GetTextMapCarrierFromMetaData(ctx context.Context) propagation.MapCarrier {
	metadataMap := make(propagation.MapCarrier)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key := range md.Copy() {
			metadataMap.Set(key, md.Get(key)[0])
		}
	}
	return metadataMap
}

func StartGrpcServerTracerSpan(ctx context.Context, operationName string) (context.Context, trace.Span) {
	textMapCarrierFromMetaData := GetTextMapCarrierFromMetaData(ctx)

	ctx = otel.GetTextMapPropagator().Extract(ctx, textMapCarrierFromMetaData)

	ctx, serverSpan := otel.GetTracerProvider().Tracer("").Start(ctx, operationName)

	return ctx, serverSpan
}

// StartKafkaConsumerTracerSpan extracts the producer's trace context from the Kafka
// message headers and starts a child span in the SAME trace (same traceId as the
// producer's request). This keeps the whole gateway -> Kafka -> core flow visible as one
// trace in Tempo, and lets Loki access-logs across all three hops be found by one TRACEID.
//
// Trade-off: since Kafka decouples producer and consumer in time (fire-and-forget, retry,
// batch delay), the producer's span/trace can appear to stay "open" until this message is
// processed - sometimes well after the producer's own request has already returned. Callers
// that care about this gap should record it explicitly (see consumer_group.go, which attaches
// kafka.produced_at / kafka.consume_delay attributes using the message's produce time).
func StartKafkaConsumerTracerSpan(ctx context.Context, headers []kafka.Header, operationName string) (context.Context, trace.Span) {
	carrierFromKafkaHeaders := TextMapCarrierFromKafkaMessageHeaders(headers)

	ctx = otel.GetTextMapPropagator().Extract(ctx, carrierFromKafkaHeaders)

	ctx, serverSpan := otel.GetTracerProvider().Tracer("").Start(ctx, operationName)

	return ctx, serverSpan
}
func TextMapCarrierToKafkaMessageHeaders(textMap propagation.MapCarrier) []kafka.Header {
	headers := make([]kafka.Header, 0, len(textMap))

	for key, val := range textMap {
		headers = append(headers, kafka.Header{
			Key:   key,
			Value: []byte(val),
		})
	}

	return headers
}

func TextMapCarrierFromKafkaMessageHeaders(headers []kafka.Header) propagation.MapCarrier {
	textMap := make(map[string]string, len(headers))
	for _, header := range headers {
		textMap[header.Key] = string(header.Value)
	}
	return propagation.MapCarrier(textMap)
}

func InjectTextMapCarrier(ctx context.Context) (propagation.MapCarrier, error) {
	m := make(propagation.MapCarrier)
	otel.GetTextMapPropagator().Inject(ctx, m)
	return m, nil
}

func InjectTextMapCarrierToGrpcMetaData(ctx context.Context, spanCtx trace.SpanContext) context.Context {
	if textMapCarrier, err := InjectTextMapCarrier(ctx); err == nil {
		md := metadata.New(textMapCarrier)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func GetKafkaTracingHeadersFromSpanCtx(ctx context.Context) []kafka.Header {
	textMapCarrier, err := InjectTextMapCarrier(ctx)
	if err != nil {
		return []kafka.Header{}
	}

	kafkaMessageHeaders := TextMapCarrierToKafkaMessageHeaders(textMapCarrier)
	return kafkaMessageHeaders
}

func TraceErr(span trace.Span, description string, err error) error {
	span.RecordError(err)
	span.SetStatus(codes.Error, description)
	return err
}
