.PHONY: run test test-unit test-integration migrate lint docker-up docker-down swagger

# Start postgres + flyway, then run API locally
run: docker-up
	cd backend && go run ./cmd/api

# Run all tests
test:
	cd backend && go test ./...

# Unit tests only (exclude integration)
test-unit:
	cd backend && go test -short ./...

# Integration tests only
test-integration:
	cd backend && go test -run Integration ./...

# Run Flyway migrations via docker-compose
migrate:
	docker compose up flyway --build

# Lint Go code
lint:
	cd backend && golangci-lint run ./...

# Start infrastructure services
docker-up:
	docker compose up -d postgres
	docker compose up flyway

# Stop all services
docker-down:
	docker compose down

# Generate Swagger/OpenAPI docs
swagger:
	cd backend && swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
