# Variables
DOCKER_COMPOSE = docker-compose
DOCKER_FILE = docker-compose.yml

# Common commands
UP = up
BUILD = build
DETACHED = -d
BUILD_FLAG = --build

# Swagger variables
SWAGGER_FILE = docs/swagger.yaml
SWAGGER_UI_PORT = 80

# Build services
build:
	$(DOCKER_COMPOSE) $(BUILD)

build-plain:
	$(DOCKER_COMPOSE) $(BUILD) --progress=plain

# Run Docker only for databases
db:
	$(DOCKER_COMPOSE) -f $(DOCKER_FILE) up $(DETACHED) postgres

# Run Docker for production
prod:
	$(DOCKER_COMPOSE) -f $(DOCKER_FILE) $(UP)

prod-d:
	$(DOCKER_COMPOSE) -f $(DOCKER_FILE) $(UP) $(DETACHED)

prod-b:
	$(DOCKER_COMPOSE) -f $(DOCKER_FILE) $(UP) $(BUILD_FLAG)

prod-bd:
	$(DOCKER_COMPOSE) -f $(DOCKER_FILE) $(UP) $(BUILD_FLAG) $(DETACHED)

# Run Docker for development
dev:
	$(DOCKER_COMPOSE) -f $(DOCKER_FILE) -f docker-compose.dev.yml $(UP)

dev-d:
	$(DOCKER_COMPOSE) -f $(DOCKER_FILE) -f docker-compose.dev.yml $(UP) -$(DETACHED)

dev-b:
	$(DOCKER_COMPOSE) -f $(DOCKER_FILE) -f docker-compose.dev.yml $(UP) $(BUILD_FLAG)

dev-bd:
	$(DOCKER_COMPOSE) -f $(DOCKER_FILE) -f docker-compose.dev.yml $(UP) $(BUILD_FLAG) $(DETACHED)

# Mock commands
mocks:
	go generate ./cmd/api/app/clients/mockgen.go
	go generate ./cmd/api/app/utils/mockgen.go
	go generate ./cmd/user/app/repositories/mockgen.go
	go generate ./cmd/security/app/repositories/mockgen.go
	go generate ./cmd/transaction/app/repositories/mockgen.go
	go generate ./cmd/broker/app/repositories/mockgen.go
	go generate ./internal/password/mockgen.go
	go generate ./pkg/email/mockgen.go
	go generate ./pkg/translation/mockgen.go

# Proto commands
proto-gen:
	protoc proto/*.proto --proto_path=proto --go_out=gen/go --go-grpc_out=gen/go

# Swagger commands
swagger: swagger-init swagger-ui swagger-gen

swagger-init:
	swag init -d cmd/api,internal -ot yaml

swagger-ui:
	docker run --rm -d -p $(SWAGGER_UI_PORT):8080 -e SWAGGER_JSON=/tmp/swagger.yaml -v `pwd`/docs:/tmp swaggerapi/swagger-ui
	@echo "Swagger UI is available at http://localhost:$(SWAGGER_UI_PORT)"

swagger-gen:
	docker run --rm -v `pwd`:/local openapitools/openapi-generator-cli:v7.11.0 generate -i /local/docs/swagger.yaml -g typescript-angular -o /local/docs/angular

# Help command to display usage
help:
	@echo "Usage:"
	@echo "  make build               \- Build all Docker services"
	@echo "  make build-plain         \- Build all Docker services with plain progress"
	@echo "  make db         	      \- Run Docker only database services"
	@echo "  make prod                \- Run Docker in production mode"
	@echo "  make prod-d              \- Run Docker in production mode detached"
	@echo "  make prod-b              \- Run Docker in production mode with build"
	@echo "  make prod-bd             \- Run Docker in production mode with build and detached"
	@echo "  make dev                 \- Run Docker in development mode"
	@echo "  make dev-d               \- Run Docker in development mode detached"
	@echo "  make dev-b               \- Run Docker in development mode with build"
	@echo "  make dev-bd              \- Run Docker in development mode with build and detached"
	@echo "  make mocks           	  \- Generate mocks for the project"
	@echo "  make proto-gen           \- Generate Go code from proto files"
	@echo "  make swagger             \- Generate and serve Swagger documentation"
	@echo "  make swagger-init        \- Initialize Swagger documentation"
	@echo "  make swagger-ui          \- Serve Swagger UI"
	@echo "  make swagger-gen         \- Generate TypeScript Angular client from Swagger"