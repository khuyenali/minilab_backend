.PHONY: run build test fmt lint sqlc-generate migrate-up migrate-down migrate-force migrate-version dev dev-setup docker-up docker-down docker-logs help swagger-init swagger-generate swagger-clean

# Variables
BINARY_NAME=mini-lab-api
MIGRATIONS_PATH=./migrations
DATABASE_URL=postgres://postgres:password@localhost:5432/minilab?sslmode=disable

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## run: Run the application
run:
	go run cmd/api/main.go

## dev: Run the application with hot reload using Air
dev:
	@echo "Starting development server with Air..."
	@echo "The server will automatically reload when you make changes to .go files"
	@echo "Press Ctrl+C to stop"
	air

## dev-setup: Setup everything needed for development
dev-setup: install-tools docker-up
	@echo "Setting up development environment..."
	@echo "Waiting for database to be ready..."
	sleep 5
	make migrate-up
	@echo "Development setup complete!"
	@echo "Run 'make dev' to start development server with hot reload"
	@echo "Visit http://localhost:8081 for Adminer (database management)"

## build: Build the application
build:
	go build -o bin/$(BINARY_NAME) cmd/api/main.go

## test: Run tests
test:
	go test -v ./...

## test-cover: Run tests with coverage
test-cover:
	go test -cover -v ./...

## sqlc-generate: Generate Go code from SQL
sqlc-generate:
	sqlc generate

## swagger-init: Initialize Swagger docs (run once)
swagger-init:
	swag init -g cmd/api/main.go -o docs

## swagger-generate: Generate/update Swagger documentation
swagger-generate:
	swag init -g cmd/api/main.go -o docs

## swagger-clean: Clean generated Swagger files
swagger-clean:
	rm -rf docs/

## docker-up: Start PostgreSQL and Adminer with Docker Compose
docker-up:
	docker compose up -d
	@echo "Services started:"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Adminer:    http://localhost:8081"

## docker-down: Stop Docker Compose services
docker-down:
	docker compose down

## docker-logs: Show Docker Compose logs
docker-logs:
	docker compose logs -f

## migrate-up: Run database migrations up
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" -verbose up

## migrate-down: Run database migrations down
migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" -verbose down

## migrate-force: Force migration to specific version (use with VERSION=n)
migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(VERSION)

## migrate-version: Check current migration version
migrate-version:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version

## migrate-create: Create new migration file (use with NAME=migration_name)
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)

## install-tools: Install development tools
install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

## clean: Clean build artifacts
clean:
	rm -f bin/$(BINARY_NAME)
	rm -rf tmp/
	go clean

## setup: Complete project setup (install tools, generate docs, etc.)
setup: install-tools deps swagger-generate
	@echo "Project setup complete!"
	@echo "Run 'make run' to start the server"
	@echo "Visit http://localhost:8080/swagger/index.html for API documentation" 