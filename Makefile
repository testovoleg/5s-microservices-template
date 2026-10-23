.PHONY: run_gw run_core run_graph init_project docker_dev local down-local clean logs-local \
	tidy deps-reset deps-upgrade deps-cleancache run-linter swagger proto_kafka proto_core \
	evans_core gqlgen

run_gw:
	go run api_gateway_service/cmd/main.go -config=./api_gateway_service/config/config.yaml

run_core:
	go run core_service/cmd/main.go -config=./core_service/config/config.yaml

run_graph:
	go run graphql_service/cmd/main.go -config=./graphql_service/config/config.yaml

# ==============================================================================
# Bootstrap new project from this template
#
#   make init_project
#   → 1/2 MODULE (e.g. 5s-dataset-api)
#   → 2/2 GQL_SUFFIX — short Federation prefix instead of tmp_
#                      (as in supergraph: dts, acc, chat → dts_eventGet, …)
#   → ShortMicroserviceName = name without 5s-  (5s-dataset-api → dataset-api)
#
#   make init_project MODULE=5s-dataset-api GQL_SUFFIX=dts
# Optional flags: SKIP_GQLGEN=1 DRY_RUN=1 FORCE=1

init_project:
	@MODULE="$(MODULE)" GQL_SUFFIX="$(GQL_SUFFIX)" \
		SKIP_GQLGEN="$(SKIP_GQLGEN)" DRY_RUN="$(DRY_RUN)" FORCE="$(FORCE)" \
		bash .data/scripts/init_project.sh

# ==============================================================================
# Docker

docker_dev:
	@echo Starting local docker dev compose
	docker-compose -f docker-compose.yaml up --build

local:
	@echo Starting local docker compose
	docker-compose -f docker-compose.local.yaml up -d --build


# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq 2>/dev/null)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)


# ==============================================================================
# Modules support

tidy:
	go mod tidy

deps-reset:
	git checkout -- go.mod
	go mod tidy

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy

deps-cleancache:
	go clean -modcache


# ==============================================================================
# Linters https://golangci-lint.run/usage/install/

run-linter:
	@echo Starting linters
	golangci-lint run ./...

# ==============================================================================
# PPROF

pprof_heap:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/heap?seconds=10

pprof_cpu:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/profile?seconds=10

pprof_allocs:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/allocs?seconds=10



# ==============================================================================
# Go migrate postgresql https://github.com/golang-migrate/migrate

DB_NAME = products
DB_HOST = localhost
DB_PORT = 5432
SSL_MODE = disable
MIGRATE_PATH = .data/migrations

force_db:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path $(MIGRATE_PATH) force 1

version_db:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path $(MIGRATE_PATH) version

migrate_up:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path $(MIGRATE_PATH) up 1

migrate_down:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path $(MIGRATE_PATH) down 1

# Прогон всех миграций up → down → up на временной БД (создаётся и удаляется): ловит
# сломанный порядок DROP/DELETE в down-файлах и ошибки в up без риска для рабочей базы.
# Нужен psql в контейнере (MIGRATE_CHECK_PSQL) либо локальный psql.
MIGRATE_CHECK_DB = $(DB_NAME)_migcheck
MIGRATE_CHECK_PSQL = docker exec microservices_postgesql psql -U postgres
migrate_check:
	$(MIGRATE_CHECK_PSQL) -c "DROP DATABASE IF EXISTS $(MIGRATE_CHECK_DB)" -c "CREATE DATABASE $(MIGRATE_CHECK_DB)"
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(MIGRATE_CHECK_DB)?sslmode=$(SSL_MODE) -path $(MIGRATE_PATH) up
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(MIGRATE_CHECK_DB)?sslmode=$(SSL_MODE) -path $(MIGRATE_PATH) down -all
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(MIGRATE_CHECK_DB)?sslmode=$(SSL_MODE) -path $(MIGRATE_PATH) up
	$(MIGRATE_CHECK_PSQL) -c "DROP DATABASE $(MIGRATE_CHECK_DB)"

# ==============================================================================
# Tests

test:
	go vet ./...
	go test ./...

fmt:
	gofmt -l -w .

# ==============================================================================
# Swagger

swagger:
	@echo Starting swagger generating
	swag init -g api_gateway_service/cmd/main.go

# ==============================================================================
# Proto

proto_kafka:
	@echo Generating kafka proto
	cd proto/kafka && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. *.proto

proto_core:
	@echo Generating core microservice proto
	cd core_service/proto && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. *.proto

# ==============================================================================
# Evans

evans_core:
	@echo Run grpc client with core proto
	cd core_service/proto && evans core_svc.proto -p 5003

# ==============================================================================
# gqlgen

gqlgen:
	@echo Generating GraphQL modules
	cd .data/gqlgen && go run -mod=mod github.com/99designs/gqlgen --verbose generate