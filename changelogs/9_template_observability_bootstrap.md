# 9. Observability, bootstrap и общие GQL/proto улучшения шаблона

## TRACEID в access-логах (GW / GQL / Core / Kafka)

Чтобы Grafana Loki → Tempo работала одинаково для всех процессов:

- HTTP access-лог пишет поле `TRACEID` (GW, GQL).
- gRPC interceptor в Core стартует server-span и пишет `TRACEID` в access-лог.
- Kafka process/commit логи также содержат `TRACEID`.
- Хелпер: `tracing.TraceIDFromContext`.
- OTL provider поднимается **до** старта HTTP/gRPC серверов.

Что смотреть при переносе:

- `pkg/tracing/utils.go` — `TraceIDFromContext`
- `pkg/interceptors/manager.go` — span в Logger interceptor
- middlewares RequestLogger (GW/GQL)
- Kafka consumer logs в core

## Порядок middleware Echo

`otelecho` регистрируется первым в `mapRoutes` (outer), затем access-logger — иначе TRACEID в логах пустой.

## gRPC: interceptor vs handler spans

- Server interceptor (`interceptorManager.Logger`, зарегистрирован в `grpc_server.go` через `grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(...))`): `StartGrpcServerTracerSpan` (extract из incoming metadata + root server span). Выполняется **один раз на каждый RPC**, до вызова конкретного delivery-хендлера.
- Client interceptor (`interceptors.TracingClientInterceptor`, зарегистрирован в `grpc.WithChainUnaryInterceptor(...)` при создании соединения в `internal/client/core_service.go` каждого вызывающего сервиса): инжектит `traceparent` текущего span'а в исходящий gRPC metadata. Тоже выполняется **один раз на каждый RPC**.
- Delivery-хендлер (`grpcService.XXX` в `core_service/internal/app/delivery/grpc/*.go`) и command-хендлеры вызывающей стороны (`api_gateway_service`/`graphql_service`): только `tracing.StartSpan` (child) — **не** вызывать `StartGrpcServerTracerSpan` и **не** вызывать `tracing.InjectTextMapCarrierToGrpcMetaData` вручную. Инжект/экстракт централизованы на уровне интерцепторов намеренно: ручной вызов в каждом хендлере легко забыть при добавлении нового метода, и тогда core-спаны молча уезжают в отдельный от gateway trace (баг был воспроизведён и исправлен в сиблинг-проекте, см. `5s-vin-decoder`).

## Kafka: producer/consumer spans

- Producer: инжект `traceparent` в `kafka.Message.Headers` сделан **один раз, централизованно**, внутри `pkg/kafka/producer.go` (`PublishMessageInTopic` → `tracing.GetKafkaTracingHeadersFromSpanCtx(ctx)`). Вызывающему коду (командным хендлерам) ничего вручную делать не нужно.
- Consumer: `coreMessageProcessor.ProcessMessages` (`core_service/internal/app/delivery/kafka/consumer_group.go`) — единственное место, где вызывается `tracing.StartKafkaConsumerTracerSpan(ctx, m.Headers, ...)`, один раз на сообщение. Получившийся span в том же trace, что и у продьюсера (единый `traceId` на весь путь gateway/graphql → Kafka → core, как и для gRPC).
- Все нижестоящие обработчики топиков (`process*` в `processor.go`) и командный слой: только `tracing.StartSpan(ctx, ...)` (child) — **не** повторный `StartKafkaConsumerTracerSpan` с `m.Headers`. Повторный extract из тех же заголовков не рвёт trace id, но ставит новый span сиблингом `ProcessMessages` вместо child, ломая дерево спанов в Tempo.
- Так как Kafka разрывает синхронность во времени (producer отвечает сразу, consumer может забрать сообщение много позже), `ProcessMessages` явно проставляет на span атрибуты `kafka.produced_at` и `kafka.consume_delay` (на основе `kafka.Message.Time`) — это даёт понятный порядок событий в Tempo/Loki, даже когда общий trace растянут во времени.

## Bulk-батчи и очередь спанов

См. [10_bulk_tracing_span_overflow.md](10_bulk_tracing_span_overflow.md): `tracing.WithBulkTracing` + увеличенные лимиты `BatchSpanProcessor`. Подключать в сервисе, когда появляется цикл по большому числу items с repo-спанами внутри.

## Health отдельно от business routes

Probes на `probes.port` (`/ready`, `/live`). Маршрут `GET /v1/health` на business HTTP убран.  
`ignoreLogUrls` включает `metrics`, `/ready`, `/live`, `/health`.

## Error-path в commands

`attrErr` + `httpErrors.ErrorWrap` на error-return в core commands (params/data в span, TraceID в тексте ошибки).

## Bootstrap нового проекта

```bash
make init_project
```

Скрипт: `.data/scripts/init_project.sh`  
Меняет module path, `ShortMicroserviceName`, GraphQL prefix `tmp_*` → `<GQL_SUFFIX>_*`.

## Proto generation

`make proto_core` / `make proto_kafka` генерируют все `*.proto` в папке — новые файлы добавлять в Makefile не нужно.

Messages core разделены:

- `core_api.proto`
- `core_event.proto`
- `core_svc.proto` (только service)

## QueryAwareBinder

`QueryAwareBinder` (query-параметры и на POST) живёт в `pkg/http_utils` и подключается в API Gateway (`e.Binder`).

## Общие GQL типы шаблона

В `graphql_service/schema`: `GeneralParamsInput`, `TipicalResponse*`, `TipicalEvent*`, пагинация, `PresignURL*`.

Несовместимо для клиентов подграфа:

- `GeneralParamsInput`: поле `accessToken` убрано (токен берётся только из заголовка `Authorization`), `companyUuid` обязателен (`String!`), добавлен `apiUuid`.
- `tmp_tipicalMutation`: аргумент `input:` переименован в `params:`; ошибка приходит ошибкой GraphQL (`data: null`), а не ответом `{code: 400, success: false, message}`.
- CORS: из `AllowOrigins` убран `https://delta2.5-systems.ru`, добавлен `https://api.5systems.ru`; `AllowHeaders` — явный список вместо `*`.
