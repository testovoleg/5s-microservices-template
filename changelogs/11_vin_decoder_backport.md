# 11. Перенос общих улучшений из 5s-vin-decoder

Пакет изменений, обкатанных в `5s-vin-decoder` (нагрузка AddVin-батчами, несколько подов
core в Kubernetes). Всё ниже не привязано к предметной области VIN.

## Репозиторий

- `.gitattributes`: `*.go`, `*.sql`, `*.proto`, `*.md`, `*.yaml/yml`, `Makefile`, `*.sh` — всегда LF.
  gofmt пишет LF, а `core.autocrlf=true` на Windows иначе превращает каждый
  отформатированный файл в дифф из всех строк. После добавления в существующий сервис —
  `git add --renormalize .` одним коммитом.
- Makefile:
  - `MIGRATE_PATH = .data/migrations` — раньше таргеты `migrate_*` смотрели в несуществующий
    `migrations/`.
  - `migrate_check` — прогон `up → down -all → up` на временной БД `<DB_NAME>_migcheck`
    (создаётся/удаляется через `psql` в контейнере Postgres). Ловит сломанный порядок
    DROP в down-файлах до того, как он сломает рабочую базу.
  - `test` (`go vet` + `go test ./...`), `fmt` (`gofmt -l -w .`).
  - `FILES := $(shell docker ps -aq 2>/dev/null)` — `make` без Docker больше не сыплет ошибку.

## pkg

- `pkg/kafka`:
  - `Producer.PublishKeyedMessage(ctx, topic, key, value)` — публикация с ключом партиции.
  - `keyAwareBalancer` в writer'ах (producer и consumer-group): сообщения с ключом
    хэшируются в одну партицию (порядок обработки одной сущности — компания, префикс),
    без ключа — прежний `LeastBytes` (равномерный bulk).
- `pkg/logger`: `InfofTrace(traceID, tmpl, args...)` — Infof с полем `TRACEID`, чтобы
  строку можно было сопоставить с `KafkaProcessMessage` того же сообщения (Loki → Tempo).
- `pkg/postgres`: `Config.MaxConns` (yaml `maxConns`, env `POSTGRES_MAX_CONNS`; 0 → 15).
  В Kubernetes суммарно по подам не должен превышать `max_connections` Postgres.
- `pkg/utils/ttlcache.go`: `TTLCache[V]` — in-process кэш «ключ → значение с TTL»
  (`Get/Put/GetOrLoad/Delete/Reset`) для редко меняющихся справочников. Переполнение `max`
  сбрасывает карту целиком; ошибки `load` не кэшируются.
- `pkg/jobs/periodic.go`: `RunPeriodic(ctx, log, name, interval, jitter, startImmediately, tick)`
  — общий цикл фоновых джоб: джиттер (поды не бьют в базу синхронно), `recover` на каждом
  тике (паника одной джобы не роняет процесс), остановка по `ctx`.

## core_service

- **Graceful shutdown** (`server.go`): после `grpcServer.GracefulStop()` ждём завершения
  `ConsumeTopic` (канал `consumerDone`, предел `consumerDrainTimeout = 30s` — должен быть
  меньше `terminationGracePeriodSeconds` пода). Раньше `defer`-ы закрывали Redis/Kafka под
  ещё живыми воркерами.
- `initKafkaTopics` / `getConsumerGroupTopics` (`utils.go`): хелпер `kafkaTopic(cfg)` —
  дефолты `partitions=1`, `replicationFactor=1` (брокер отвергает `0`), пустое имя топика
  (выключен конфигом) пропускается. Новый топик — по одной строке в оба списка.
- **Readiness Kafka**: свежее соединение на каждую проверку с `probeTimeout = 5s`.
  Долгоживущий `s.kafkaConn` после переезда брокера остаётся мёртвым, и под навсегда
  выпадал из Service при живом liveness. Для Postgres (если подключён) — тот же паттерн:
  `pgConn.Ping(pingCtx)` с таймаутом.
- **Keycloak introspect** (`keycloak_repository.go`): `gocloak.RetrospectToken` жёстко шлёт
  `token_type_hint=requesting_party_token` (UMA RPT), и Keycloak отвечает `active:false` на
  обычные access-токены. Теперь `introspectAccessToken` вызывает
  `/protocol/openid-connect/token/introspect` напрямую с `token_type_hint=access_token`
  (Basic auth клиентом сервиса). При `active:false` в лог уходят iss/azp/sid/exp токена
  (без проверки подписи, только для диагностики «не тот realm/клиент»).
- `adminRepository.GetCompanyList` — весь справочник компаний Admin API. Пагинация
  `/company/list` **1-based** (обход с 0 дублирует первую страницу), предел
  `AdminCompanyListMaxPages`. Кэш `CacheRepository.GetCompanyList/PutCompanyList`
  (TTL 5 мин) — для админских выборок, не для горячего пути.

## api_gateway_service

- `http.bodyLimit` (yaml, например `20M`; пусто → `2M`) вместо константы — bulk-импорты
  не влезали в 2M.
- **Своя index-страница Swagger UI** (`delivery/http/v1/swagger.go` + `swagger_ui.html`
  через `//go:embed`): `GET /swagger/*` отдаёт её вместо stock-страницы, остальное (`doc.json`,
  js, css) — echo-swagger. Заголовок — из `http.title`.
  - Плагин UI: для тел > 50 000 символов не строится curl-сниппет (swagger-ui падал с
    «Maximum call stack size exceeded» на больших JSON).
  - `swagger_test.go` рендерит шаблон: опечатка в нём иначе всплыла бы паникой
    `template.Must` при первом заходе на Swagger.

## Что смотреть при переносе в существующий сервис

1. `.gitattributes`, Makefile (`MIGRATE_PATH`, `migrate_check`, `test`, `fmt`).
2. `pkg/kafka/{producer,writer,consumer_group}.go`, `pkg/logger/logger.go`,
   `pkg/postgres/postgres.go`, `pkg/constants/constants.go`, `pkg/utils/ttlcache.go`, `pkg/jobs/`.
3. `core_service/internal/server/{server,utils}.go` — drain воркеров, `kafkaTopic`, readiness.
4. `core_service/internal/app/repository/keycloak_repository.go` — introspect.
5. `api_gateway_service`: `config/config.{go,yaml}` (`http.bodyLimit`),
   `internal/server/http.go`, `delivery/http/v1/{routes,swagger,swagger_test}.go`,
   `swagger_ui.html`.
