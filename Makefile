include .env
export

DB_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable

.PHONY: help run up down migration migrate-up migrate-down

help:
	@echo Available commands:
	@echo   - make up                    Start PostgreSQL
	@echo   - make down                  Stop PostgreSQL
	@echo   - make run                   Run application
	@echo   - make migration name=NAME   Create migration
	@echo   - make migrate-up            Apply migrations
	@echo   - make migrate-down          Rollback last migration

up:
	docker compose up -d

down:
	docker compose down

run:
	go run ./cmd/url_shortener

migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1