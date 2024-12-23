# Добавление клиента в репозиторий
В core_service\internal\app\repository\keycloak_repository.go ( файл репозитория может называться иначе ) в переменной idmRepository добавить следующую переменную:
```go
client *gocloak.GoCloak
```
В функции NewIDMRepository добавить переменную клиента. К примеру получиться скорее всего так:
```go
func NewIDMRepository(log logger.Logger, cfg *config.Config, keycloakClient *gocloak.GoCloak) *idmRepository {
	return &idmRepository{log: log, cfg: cfg, client: keycloakClient}
}
```
В методах репозитория убрать создание клиента и проставить чтобы клиент тянулся из переменной репозитория. Пример:
```go
p.client.RetrospectToken()
```
# Добавление клиента keycloak в сервер
В core_service\internal\server\server.go в переменную сервер добавить
```go
keycloakClient *gocloak.GoCloak
```
В функции Run() добавить следующее перед иннициализацией репозиториев
```go
s.keycloakClient = gocloak.NewClient(s.cfg.Keycloak.Host)
```
Добавить новую переменную при инниацилизации репозитория
```go
cloakRepo := repository.NewIDMRepository(s.log, s.cfg, s.keycloakClient)
```
