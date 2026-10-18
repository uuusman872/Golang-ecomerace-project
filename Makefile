.PHONY: help build run dev lint migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  build         Build the application binary"
	@echo "  run           Run the application"
	@echo "  dev           Run the application in development mode"
	@echo "  lint          Run golangci-lint over the codebase"
	@echo "  migrate-up    Apply all up database migrations"
	@echo "  migrate-down  Revert all down database migrations"

build:
	go build -o bin/app ./cmd/api

run:
	go run ./cmd/api

dev:
	go run ./cmd/api

lint:
	golangci-lint run ./...

migrate-up:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable" down

docker-up:
	sudo docker compose -f docker/docker-compose.yml up -d

docker-down:
	sudo docker compose -f docker/docker-compose.yml down