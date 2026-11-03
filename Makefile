.PHONY: help build run dev lint migrate-up migrate-down swagger

help:
	@echo "Available commands:"
	@echo "  build         Build the application binary"
	@echo "  run           Run the application"
	@echo "  dev           Run the application in development mode"
	@echo "  lint          Run golangci-lint over the codebase"
	@echo "  format        Run gofmt and goimports to format the code"
	@echo "  migrate-up    Apply all up database migrations"
	@echo "  migrate-down  Revert all down database migrations"
	@echo "  swagger       Generate Swagger docs (swag init)"
	@echo "  docker-up     Build and start full stack (postgres, api, localstack, mailpit, cdn)"
	@echo "  docker-down   Stop docker compose stack"
	@echo "  docker-logs   Follow api container logs"

build:
	@echo "Building all binaries...."
	@mkdir -p bin
	@for cmd in cmd/*/; do \
    		if [ -d "$$cmd" ]; then \
    			binary=$$(basename $$cmd); \
    			echo "Building $$binary..."; \
    			go build -o bin/$$binary ./$$cmd; \
    		fi \
    	done

run:
	go run ./cmd/api

dev:
	go run ./cmd/api

lint:
	golangci-lint run ./...

format:
	@gofmt -s -w .
	@goimports -w .

migrate-up:
	migrate -path db/migrations -database "postgresql://root:admin1234@localhost:5432/ecomerace_shop?sslmode=disable" up  

migrate-down:
	migrate -path db/migrations -database "postgresql://root:admin1234@localhost:5432/ecomerace_shop?sslmode=disable" down

docker-up:
	docker compose --env-file .env -f docker/docker-compose.yml up -d --build

docker-down:
	docker compose --env-file .env -f docker/docker-compose.yml down

docker-logs:
	docker compose --env-file .env -f docker/docker-compose.yml logs -f api

swagger:
	swag init -g cmd/api/main.go -o docs --parseInternal

graph-generate:
	@go get github.com/99designs/gqlgen@v0.17.78
	@go run github.com/99designs/gqlgen generate