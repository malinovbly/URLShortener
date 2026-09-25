include .env
export

.PHONY: help swagger up down build run migration migrate-up migrate-down

help:
	@echo "Available commands:"
	@echo "  make swagger               Generate Swagger (OpenAPI)"
	@echo "  make up                    Start all containers"
	@echo "  make down                  Stop all containers"
	@echo "  make build                 Build application image"
	@echo "  make run                   Build and start all containers"
	@echo "  make migration name=NAME   Create migration"
	@echo "  make migrate-up            Apply migrations"
	@echo "  make migrate-down          Rollback last migration"

swagger:
	swag init -g cmd/url_shortener/main.go

up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

run:
	docker compose up --build -d

migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	docker compose run --rm migrate

migrate-down:
	docker compose run --rm migrate \
		-path=/migrations \
		-database="postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable" \
		down 1