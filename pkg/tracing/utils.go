package tracing

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.GetTracerProvider().Tracer("").Start(ctx, name)
}

func StartHttpServerTracerSpan(c echo.Context, operationName string) (context.Context, trace.Span) {
	ctx := c.Request().Context()

	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request().Header))

	ctx, serverSpan := otel.GetTracerProvider().Tracer("").Start(ctx, operationName)

	traceID := serverSpan.SpanContext().TraceID()
	if traceID.IsValid() {
		c.Set("traceid", traceID.String())
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

func StartKafkaConsumerTracerSpan(ctx context.Context, headers []kafka.Header, operationName string) (context.Context, trace.Span) {
	carrierFromKafkaHeaders := TextMapCarrierFromKafkaMessageHeaders(headers)

	ctx = otel.GetTextMapPropagator().Extract(ctx, carrierFromKafkaHeaders)

	ctx, serverSpan := otel.GetTracerProvider().Tracer("").Start(ctx, operationName)
	//ctx = trace.ContextWithSpan(ctx, spanCtx)

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
