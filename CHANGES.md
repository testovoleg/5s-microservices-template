
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



