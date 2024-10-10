
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


