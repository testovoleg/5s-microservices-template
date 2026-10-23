# 5s microservices template

Корпоративный шаблон микросервиса 5Systems (Go). Один `go.mod`, три процесса, общий `pkg/`.

## Архитектура

```text
Client
  ├─ REST     → api_gateway (:5001) ──gRPC──► core (:5002)
  ├─ GraphQL  → graphql (:5010)     ──gRPC──► core
  └─ Webhook  → gateway ──Kafka──► core consumer
                                      ├─ Keycloak / Redis
                                      └─ внешние 5s API (admin/auth/storage)
```

| Сервис | Роль | Порт |
|--------|------|------|
| `api_gateway_service` | REST + Swagger/Redoc, Kafka producer | `:5001` |
| `core_service` | Бизнес-логика: gRPC + Kafka consumers | `:5002` |
| `graphql_service` | GraphQL BFF (federation subgraph) | `:5010` |

Kubernetes probes (`/ready`, `/live`) слушают на отдельном порту `probes.port`, не на бизнес-HTTP: GW `:3001`, Core `:3003`, GQL `:3010`.  
pprof — `probes.pprof` (GW `:6001`, Core `:6003`, GQL `:6010`). Metrics — `probes.prometheusPort`.

## Стек

- Go 1.25, Echo, gRPC, gqlgen (federation v2)
- Kafka, Redis, Keycloak
- OpenTelemetry → OTLP (`otl.endpoint`)
- Опционально Postgres (`pkg/postgres`) — в шаблоне конфиг есть, подключение в core не включено; используется в проектах вроде `5s-data-coordinator`
- Docker-образы сервисов лежат в `.docker/` (для сборки/деплоя). Локальный `docker-compose` в шаблоне устарел и не является основным способом разработки.

## Быстрый старт (локально)

Нужны: Go, Redis, Kafka, (опционально) Keycloak / OTL-collector.

```bash
# правки конфигов под своё окружение
# */config/config.yaml

make run_core    # gRPC :5002
make run_gw      # REST  :5001
make run_graph   # GQL   :5010
```

Полезные make-таргеты:

| Target | Назначение |
|--------|------------|
| `init_project` | Bootstrap нового репозитория из шаблона |
| `run_core` / `run_gw` / `run_graph` | Запуск процессов |
| `proto_core` / `proto_kafka` | Генерация `*.proto` → Go (`*.proto` в папке) |
| `gqlgen` | Генерация GraphQL |
| `swagger` | Swagger для API Gateway |
| `tidy` / `run-linter` | Зависимости / линтер |
| `test` / `fmt` | `go vet` + `go test ./...` / `gofmt -w` |
| `migrate_up` / `migrate_down` / `migrate_check` | Миграции `.data/migrations`; `migrate_check` — up → down → up на временной БД |

Swagger UI: `http://localhost:5001/v1/swagger/index.html`  
GraphQL playground: `http://localhost:5010/graphql/playground`

## Создание нового сервиса из шаблона

1. Скопировать репозиторий / создать ветку от шаблона.
2. `make init_project`
   - **MODULE** — имя проекта, например `5s-dataset-api` → module `github.com/testovoleg/5s-dataset-api`
   - **GQL_SUFFIX** — короткий префикс полей в federation (как в supergraph: `dts`, `acc`, `chat` → `dts_eventGet`)
   - `ShortMicroserviceName` выставляется автоматически (`5s-dataset-api` → `dataset-api`) — префикс Redis-ключей
3. Проверить diff: `go.mod`, импорты, `pkg/constants`, GraphQL `tmp_*` → `<suffix>_*`.
4. Подправить `*/config/config.yaml` (порты, kafka topics, OTL, Keycloak, внешние API).
5. Заменить tipical/example domain-stubs на реальные команды/репозитории.
6. Регистрацию subgraph в `5s-supergraph` делать вручную в том репозитории (сюда не входит).

Зачем этот чеклист: `init_project` меняет имена и импорты, но не знает порты/топики/внешние URL конкретного сервиса. Без явных шагов 4–6 сервис «собирается», но не встаёт в окружение.

## Observability

- Access-логи GW / GQL / Core (gRPC) / Kafka содержат поле `TRACEID` для связки Loki → Tempo.
- HTTP-ответы могут отдавать заголовок `X-Trace-ID`.
- Ошибки из commands/репозиториев дополняются `TraceID:` через `pkg/http_errors`.
- gRPC client interceptor инжектит `traceparent` на каждом вызове; Kafka consumer продолжает тот же `traceId`.
- Большие батч-циклы: `tracing.WithBulkTracing(ctx)`, см. [changelogs/10_bulk_tracing_span_overflow.md](changelogs/10_bulk_tracing_span_overflow.md).

## Фоновые джобы и кэш

- `jobs.RunPeriodic` (`pkg/jobs`) — цикл фоновой джобы с джиттером и `recover` на каждом тике; запускать `go jobs.RunPeriodic(ctx, ...)` из `server.Run` после подключения репозиториев.
- `utils.TTLCache[V]` (`pkg/utils/ttlcache.go`) — in-process кэш справочников с TTL.
- Kafka: `Producer.PublishKeyedMessage` — сообщения одной сущности в одну партицию (порядок), без ключа — равномерный `LeastBytes`.
- Выключение core: gRPC `GracefulStop` → ожидание Kafka-воркеров (`consumerDrainTimeout`, 30s) → закрытие пулов. `terminationGracePeriodSeconds` пода должен быть больше.

## Перенос улучшений в существующие сервисы

См. [CHANGES.md](CHANGES.md) и файлы в `changelogs/`.
