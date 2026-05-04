.PHONY: dev build docker-up docker-down migrate lint

# Backend dev server (requires DB running)
dev:
	go run ./cmd/app/main.go

# Build binary
build:
	go build -o bin/retail-api ./cmd/app/main.go

# Run all tests
test:
	go test ./... -v

# Docker Compose
docker-up:
	docker compose up -d

docker-down:
	docker compose down

# Run migration manually
migrate:
	docker compose exec postgres psql -U retail_user -d retail_db -f /docker-entrypoint-initdb.d/001_init.up.sql

# Copy .env
env:
	cp .env.example .env
	@echo "✅ .env created — remember to update JWT_SECRET"

# Format
fmt:
	go fmt ./...

# Lint
lint:
	go vet ./...
