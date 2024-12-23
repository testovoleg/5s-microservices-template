
# checklist

- [ ] Добавлен файл pkg/tracing/opentelemetry.go
- [ ] Удален файл pkg/tracing/jager.go
- [ ] Модифицирован файл pkg/tracing/utils.go
- [ ] Константа jaeger заменена на OTL
- [ ] Замена конфига jaeger в config/config.go на конфиг OTL
- [ ] Заменена загрузка переменной окружения среды в config/config.go ( с jaeger на OTL )
- [ ] Заменены значения по умолчанию в config/config.yaml ( убран jaeger, добавлен OTL )
- [ ] Замена провайдера трассировки в internal/server/server.go ( с jaeger на OTL )
- [ ] Убран opentracing обработчика в сервере grpc
- [ ] Изменена обработка ошибок traceErr в api_gateway
- [ ] go mod tidy
- [ ] Массовая замена
- [ ] Отсутствие ошибок при компиляции
- [ ] go mod tidy
- [ ] Установка переменных окружения в kubernetes


# 1. Добавление файла новой библиотеки в pkg
## Добавление файла
Добавление файла pkg/tracing/opentelemetry.go со следующим содержимым:
```go
package tracing

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type OTLConfig struct {
	ServiceName string `mapstructure:"serviceName"`
	Endpoint    string `mapstructure:"endpoint"`
	Enable      bool   `mapstructure:"enable"`
}

func NewOTLTracer(ctx context.Context, conf *OTLConfig) (provider *trace.TracerProvider, shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	resource, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(conf.ServiceName),
	))

	// Set up trace provider.
	tracerProvider, err := newTraceProvider(ctx, conf.Endpoint, resource)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)

	return tracerProvider, shutdown, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context, endpoint string, res *resource.Resource) (*trace.TracerProvider, error) {
	// create otlp grpc exporter
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	// create provider with exporter and custom resources
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithBatchTimeout(time.Second),
		),
		trace.WithResource(res),
	)
	return traceProvider, nil
}
```

## Удаления файла
Если есть файл pkg/tracing/jager.go , то удалить его

## Изменения utils.go
Удалить все содержимое в utils.go и заменить на следующе более актуальное ( по хорошему возьмите из проекта template ), или замените на следующее ( возможно менее актуальное ):
```go
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
```

# 2. Добавление конфига во все контейнеры

## Добавление константы
Удалить старую константу (JaegerHostPort) и добавить новую (OTLEndpoint) в pkg/constants/constants.go
```go
OTLEndpoint = "OTL_ENDPOINT"
```
## Добавление конфига в контейнеры

В core_service, api_gateway, graphql в config/config.go удалить jaeger из структуры Config
```go
Jaeger          *tracing.Config     `mapstructure:"jaeger"`
```
и добавить otl
```go
OTL             *tracing.OTLConfig  `mapstructure:"otl"`
```
В фукнции InitConfig убрать загрузку jaeger из переменных окружения ( название функции может отличатся )
```go
utils.CheckEnvStr(&cfg.Jaeger.HostPort, constants.JaegerHostPort)
```
и добавить otl загрузку переменной окружения
```go
utils.CheckEnvStr(&cfg.OTL.Endpoint, constants.OTLEndpoint)
```


## Добавление значений по умолчанию в конфиге

В core_service, api_gateway, graphql в config/config.go удалить конфиг по умолчанию jaeger
```yaml
jaeger:
  enable: true
  serviceName: 
  hostPort: "localhost:6831"
  logSpans: false
```
и добавить следующее ( подставив service_name )
```yaml
otl:
  enable: true
  serviceName: 
  endpoint: "localhost:4317"
```

# Добавление в сервере провайдера трассировки
В core_service, api_gateway, graphql в internal/server/server.go удалить создание jaeger трассировщика
```go
if s.cfg.Jaeger.Enable {
    tracer, closer, err := tracing.NewJaegerTracer(s.cfg.Jaeger)
    if err != nil {
        return err
    }
    defer closer.Close() // nolint: errcheck
    opentracing.SetGlobalTracer(tracer)
}
```
и добавить otl трассировщика
```go
if s.cfg.OTL.Enable {
    provider, shutdown, err := tracing.NewOTLTracer(ctx, s.cfg.OTL)
    if err != nil {
        return err
    }
    defer func() { _ = shutdown(ctx) }()

    otel.SetTracerProvider(provider)
}
```
## Убирание opentracing обработчика grpc
В файле core_service/internal/server/grpc_server.go закомментировать следующую строку
```go
grpc_opentracing.UnaryServerInterceptor(),
```

# Изменение обработки ошибок в api_gateway
В api_gateway_service\internal\app\delivery\http\v1\handlers.go понадобится заменить функцию "func (h *appHandlers) traceErr(span opentracing.Span, err error)" \
А именно, удалить это:
```go
func (h *appHandlers) traceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
	h.metrics.ErrorHttpRequests.Inc()
}
```
И добавить вместо этого кода это:
```go
func (h *appHandlers) traceErr(span trace.Span, err error) {
	span.SetStatus(codes.Error, "operation in api_gateway failed")
	span.RecordError(err)
	h.metrics.ErrorHttpRequests.Inc()
}
```

# Добавление библиотеки в go.mod
Пробить следующую команду в командной строке, для добавления библиотек в проект
```
go mod tidy
```

После подобной команды d go.mod должна пропасть библиотека jaeger-client-go


# Массовая замена
Теперь остается заменить все функции старой библиотеки на новую. Делать это будет по одной строке. Так же, в одном из случаев придется воспользоваться regexp заменой. \
Массовая замена будет происходить во всем проекте, а не в конкретных файлах. Примеры сделаны в рамках замены в vs code

## Первая замена
Здесь будет применяться regexp замена. Убедитесь что она у вас включена ( значок выглядит как символы ".*" ). \
Включить её нужно будет только для этой замены, после желательно выключить \
После замены проверьте что корректно проведена заменна ( сохранились названия трассировок ) \
\
Заменяем
```go
span, (ctx|_?) (:?)= opentracing.StartSpanFromContext\(ctx, "(.+)"\)
```
на 
```go
$1, span $2= tracing.StartSpan(ctx, "$3")
```
## Вторая замена
Заменяем
```go
defer span.Finish()
```
на
```go
defer span.End()
```
## Третья замена
Заменяем
```go
span.Context()
```
на
```go
span.SpanContext()
```
## Четвертая замена
Заменяем
```go
if opentracing.SpanFromContext(ctx) != nil {
```
на
```go
if trace.SpanContextFromContext(ctx).IsValid() {
```
## Пятая замена
Заменяем
```go
var span opentracing.Span
```
на
```go
var span trace.Span
```
## Шестая замена
Заменяем
```go
tracing.GetKafkaTracingHeadersFromSpanCtx(span.SpanContext())
```
на
```go
tracing.GetKafkaTracingHeadersFromSpanCtx(ctx)
```

# Проверка на ошибки
После этого, ошибок в коде быть не должно, максимум предупреждение о том что у вас не используется библиотека opentracing-go \
Ошибки могут быть в случаи если были какие то трассировки сделаны как то по другому, поэтому вам понадобиться самостоятельно заменить остальные штуки. \
Возможно может появится ошибка о том что не може загрузить библиотеку google.golang...cloudtrace/v2 , в таком случаи нужно просто снести загрузку этой библиотеки в файле, после автоматом должна подтянутся нормальная библиотека.

# Еще раз go mod tidy
В командной строке еще раз пробиваем команду go mod tidy для того чтобы избавится полностью от opentracing-go

# Пробивание в kubernetes переменных трассировки
Т.к. мы изменили названия переменных для трассировки, то в манифестах кубернетеса понадобиться изменить эти названия, и указывать на актуальный сервер трассировки. \
Возможно сначала стоит пробить новую переменную, а старую удалить после того как применится коммит с новой библиотекой трассировки.

