
# 1. Отложенный краш GRPC
В случаи есл

## Причины изменений 
Более детальные причины краша

## Что изменять
Обработку ошибок core_service\internal\server\grpc_server.go

## Изменения
Описаны в [changelogs/1_grpc_handle_error.md](changelogs/1_grpc_handle_error.md)

# 2. Изменение версии golang в конйтенерах
Раз в какой то период, необходимо обновлять версию golang. \
Чтобы обновить в проде версию golang, необходимо обновить эту версию в dockerfile. \

## Причины изменений
У некоторых библиотек обязательное требование golang не ниже определенной версии.

## Что изменять
Dockerfile каждого контейнера

## Изменения
Описаны в [changelogs/2_golang_change_version.md](changelogs/2_golang_change_version.md)

# 3. Изменение библиотеки трассировки
Изменяется библиотека трассировки с opentracing на opentelemetry.

## Причины изменений 
Библиотека opentracing теперь deprecated ( более не поддерживается и не обновляется )

## Что изменять
Все что связано с opentracing ( начало/конец трассировок, создание сервера трассировки и т.д. )

## Изменения
Описаны в [changelogs/3_tracing_library_change.md](changelogs/3_tracing_library_change.md)


# 4. Трассировка в репозиториях только при наличии запущенной трассировки
Если вызывать метод старта трассировки, без родительской трассировки то будет стартовать новая трассировка. \
В некоторых проектах есть фоновые задания, которые вызываются без трассировок, но пользуются методами репозиториев. Из за чего при вызове функции из репозитория каждый раз создается отдельная трассировка.

## Причины изменений 
Избавление от лишних бессмысленные трассировки

## Что изменять
Все трассировки в репозиториях ( core_service/internal/app/repository/* ).

## Изменения
Описаны в Описаны в [changelogs/4_only_started_trace.md](changelogs/4_only_started_trace.md)

# 5. Добавление traceID в логах
Данная штука позволяет добавить индетификатор трассировки в api_gateway и graphql, для того чтобы по логу найти трассировку 

## Причина изменений
Более удобный поиск соотношений лога к трассировке ( и наоборот )

# Изменения
Описаны в [changelogs/5_add_traceid_in_log.md](changelogs/5_add_traceid_in_log.md)

# 6. Вынос keycloak клиента в переменную сервера
Изначально старые версии keycloak репозитория сделаны так что в каждом методе создают собственного клиента.\
Для оптимизации, нужно создать этот клиент один раз и использовать его во всех клиентах.

## Причина изменений
Улучшение быстродействия, за счёт того что не будет создаваться новый клиент каждый запрос

## Что изменять
Методы в репозитории keycloak ( отсутствует в template )

## Изменения
Описаны в [changelogs/6_keycloak_one_create_client.md](changelogs/6_keycloak_one_create_client.md)

# 7. Добавление сквозной трассировки
Трассировки изначально не передаются между микросервисами ( если быть точнее то по http не передаются ). Для того чтобы реализовать подобное, понадобиться в header передать идентификаторы трассировки.

## Что изменять
Все методы которые используют http обращения к нашим внутренним микросервисам

## Измененния 
Описаны в [changelogs/7_traceid_passthrough.md](changelogs/7_traceid_passthrough.md)

# 8. Добавлените идентификатора трассировки в ответах api gateway
Для большего удобства дебага, было предумано сделать возврат идентификатора трассировки в api_gw. Возвращается он в header.

## Что изменять
Добавить строчку кода в pkg/tracing/utils.go

## Измененния 
Описаны в [changelogs/8_x_traceid_gw.md](changelogs/8_x_traceid_gw.md)

# 9. Observability / bootstrap / proto / GQL commons в шаблоне
Пакет улучшений шаблона: TRACEID в access-логах всех процессов (для Loki→Tempo), порядок Echo middleware, gRPC child-spans, health на probes.port, `attrErr` в commands, `make init_project`, proto через `*.proto`, общие GQL input/types.

## Причины изменений
Единый DX при создании сервиса из шаблона и одинаковая наблюдаемость GW/GQL/Core без ручной донастройки TRACEID.

## Что изменять
См. список файлов и шагов переноса в changelog.

## Изменения
Описаны в [changelogs/9_template_observability_bootstrap.md](changelogs/9_template_observability_bootstrap.md)

# 10. Bulk-батчи и переполнение очереди трассировки
На большом батче SDK молча дропает спаны из очереди `BatchSpanProcessor`; чаще всего пропадает
родитель цикла, а дети остаются «сиротами» в Tempo. В шаблоне увеличены лимиты очереди и
добавлен `tracing.WithBulkTracing` для будущих bulk-циклов.

## Причины изменений
Единое дерево parent/child в Tempo даже при обработке больших батчей.

## Что изменять
`pkg/tracing/opentelemetry.go`, `pkg/tracing/utils.go`. В конкретном сервисе — обернуть bulk-цикл
в `tracing.WithBulkTracing`.

## Изменения
Описаны в [changelogs/10_bulk_tracing_span_overflow.md](changelogs/10_bulk_tracing_span_overflow.md)

# 11. Перенос общих улучшений из 5s-vin-decoder
`.gitattributes` (LF), Makefile (`MIGRATE_PATH`, `migrate_check`, `test`, `fmt`), Kafka keyed
publish + `keyAwareBalancer`, `logger.InfofTrace`, `postgres.MaxConns`, `utils.TTLCache`,
`jobs.RunPeriodic`, graceful drain Kafka-воркеров, `kafkaTopic` с дефолтами, readiness Kafka со
свежим соединением, фикс Keycloak introspect (gocloak шлёт не тот `token_type_hint`),
`adminRepository.GetCompanyList`, `http.bodyLimit`, Swagger UI без curl-сниппета для больших тел.

## Причины изменений
Проблемы, пойманные под нагрузкой и в Kubernetes на `5s-vin-decoder`: обрыв пулов под живыми
воркерами при выключении пода, «залипший» readiness после переезда брокера, `active:false` на
валидные токены, 2M-лимит тела на bulk-импортах; плюс DX — Swagger, не падающий на больших
телах, и проверка миграций.

## Что изменять
См. список файлов в changelog.

## Изменения
Описаны в [changelogs/11_vin_decoder_backport.md](changelogs/11_vin_decoder_backport.md)
