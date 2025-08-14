.PHONY: build run test clean deps fmt lint help

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
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  compose-up   - Start all services with docker-compose"
	@echo "  compose-down - Stop all services"
	@echo "  compose-logs - View logs from all services"
	@echo "  help         - Show this help message"

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
