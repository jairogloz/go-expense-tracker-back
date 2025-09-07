.PHONY: build run test clean deps fmt lint help migrate-up migrate-up-one migrate-down-one migrate-status migrate-create

# Build the application
build:
	go build -o bin/expense-tracker cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Show help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Download and organize dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code (requires golangci-lint)"
	@echo "  db-up        - Start PostgreSQL database with Docker"
	@echo "  db-down      - Stop PostgreSQL database"
	@echo "  db-connect   - Connect to PostgreSQL database"
	@echo "  migrate-up   - Run all database migrations up"
	@echo "  migrate-up-one - Run one database migration up"
	@echo "  migrate-down-one - Run one database migration down"
	@echo "  migrate-status - Show current migration status"
	@echo "  migrate-create - Create new migration (requires name=migration_name)"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  compose-up   - Start all services with docker-compose"
	@echo "  compose-down - Stop all services"
	@echo "  compose-logs - View logs from all services"
	@echo "  help         - Show this help message"

# Database configuration (can be overridden with environment variables)
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= expense_user
DB_PASSWORD ?= expense_password
DB_NAME ?= expense_tracker
DB_SSL_MODE ?= disable

# Construct database URL
DB_URL = postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

# Database migrations (requires golang-migrate: brew install golang-migrate)
migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-up-one:
	migrate -path migrations -database "$(DB_URL)" -verbose up 1

migrate-down-one:
	migrate -path migrations -database "$(DB_URL)" -verbose down 1

migrate-status:
	migrate -path migrations -database "$(DB_URL)" version

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Development database setup (requires docker)
db-up:
	docker run --name expense-tracker-db \
		-e POSTGRES_DB=expense_tracker \
		-e POSTGRES_USER=expense_user \
		-e POSTGRES_PASSWORD=expense_password \
		-p 5432:5432 \
		-d postgres:15

# Stop development database
db-down:
	docker stop expense-tracker-db
	docker rm expense-tracker-db

# Connect to development database
db-connect:
	docker exec -it expense-tracker-db psql -U expense_user -d expense_tracker

# Docker commands
docker-build:
	docker build -t expense-tracker .

docker-run:
	docker run --rm -p 8080:8080 --env-file .env expense-tracker

# Start all services with docker-compose
compose-up:
	docker-compose up -d

# Stop all services
compose-down:
	docker-compose down

# View logs
compose-logs:
	docker-compose logs -f
