.PHONY: help up down restart logs build run-go run-ai test clean

help:
	@echo "Available commands:"
	@echo "  make up          - Start all docker containers (Postgres, Redis, RabbitMQ, Services)"
	@echo "  make down        - Stop all docker containers"
	@echo "  make restart     - Restart all docker containers"
	@echo "  make logs        - Tail logs of all containers"
	@echo "  make build       - Build Go binary and AI Docker images"
	@echo "  make run-go      - Run Go core service locally"
	@echo "  make run-ai      - Run FastAPI AI service locally"
	@echo "  make test        - Run tests for Go and FastAPI"
	@echo "  make clean       - Remove temporary files and build artifacts"

up:
	docker-compose up -d

down:
	docker-compose down

restart:
	docker-compose restart

logs:
	docker-compose logs -f

build:
	cd services/core-go && go build -o ../../bin/core_server.exe ./cmd/server/main.go

run-go:
	cd services/core-go && go run ./cmd/server/main.go

run-ai:
	cd services/ai-fastapi && uvicorn cmd.api.main:app --reload --port 8000

test:
	cd services/core-go && go test -v ./...
	cd services/ai-fastapi && pytest

clean:
	rm -rf bin/
