.PHONY: help build run test clean docker-build docker-up docker-down migrate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/auth-service cmd/server/main.go

run: ## Run the application locally
	go run cmd/server/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/

docker-build: ## Build Docker image
	docker-compose build

docker-up: ## Start all services with Docker Compose
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

migrate: ## Run database migrations (requires psql)
	psql -U postgres -d auth_db -f migrations/001_init.sql

deps: ## Download dependencies
	go mod tidy
	go mod download

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

.DEFAULT_GOAL := help
