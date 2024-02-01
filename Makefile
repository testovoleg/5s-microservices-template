.PHONY:

run_gw:
	go run api_gateway_service/cmd/main.go -config=./api_gateway_service/config/config.yaml

run_core:
	go run core_service/cmd/main.go -config=./core_service/config/config.yaml

run_graph:
	go run graphql_service/cmd/main.go -config=./graphql_service/config/config.yaml

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

FILES := $(shell docker ps -aq)

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

force_db:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path migrations force 1

version_db:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path migrations version

migrate_up:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path migrations up 1

migrate_down:
	migrate -database postgres://postgres:postgres@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE) -path migrations down 1

# ==============================================================================
# Swagger

swagger:
	@echo Starting swagger generating
	swag init -g api_gateway_service/cmd/main.go

# ==============================================================================
# Proto

proto_kafka:
	@echo Generating kafka proto
	cd proto && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. kafka.proto

proto_auth:
	@echo Generating authorization microservice proto
	cd auth_service/proto && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. --proto_path=. auth_service.proto auth_dto.proto
		
proto_core:
	@echo Generating connector microservice proto
	cd core_service/proto && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. core_svc.proto core_messages.proto

# ==============================================================================
# Evans

evans_core:
	@echo Run grpc client with core proto
	cd core_service/proto && evans core_svc.proto -p 5003

# ==============================================================================
# gqlgen

gqlgen:
	@echo Generating GraphQL modules
	cd .data/gqlgen && go run -mod=mod github.com/99designs/gqlgen --verbose gener