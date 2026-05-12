SHELL := /bin/sh

.PHONY: help backend-run backend-test backend-fmt backend-build docker-up docker-down migrate-up migrate-down

help:
	@echo "TerraRoute commands:"
	@echo "  make backend-run    Run the backend locally"
	@echo "  make backend-test   Run backend tests"
	@echo "  make backend-fmt    Format backend Go code"
	@echo "  make backend-build  Build backend binary"
	@echo "  make docker-up      Start local services"
	@echo "  make docker-down    Stop local services"
	@echo "  make migrate-up     Apply database migrations"
	@echo "  make migrate-down   Roll back database migrations"

backend-run:
	cd backend && go run ./cmd/api

backend-test:
	cd backend && go test ./...

backend-fmt:
	cd backend && go fmt ./...

backend-build:
	cd backend && go build -o bin/terraroute-api ./cmd/api

docker-up:
	docker compose up --build

docker-down:
	docker compose down

migrate-up:
	docker compose run --rm migrate

migrate-down:
	docker compose run --rm migrate -path=/migrations -database "$${API_DATABASE_URL:-postgres://terraroute:terraroute_dev_password@postgres:5432/terraroute?sslmode=disable}" down 1
