# URLShortener


## Table of Contents

- [Description](#description)
- [Running](#running)
- [Migrations](#migrations)


## Description

A URL shortening service written in Go. Uses PostgreSQL for storing URLs.

The service allows users to create short URLs and redirect users using an alias.


## Running

Start PostgreSQL:
```bash
docker compose up -d
```

Start the application:
```bash
go run ./cmd/url_shortener
```


## Migrations

Create a migration:
```bash
migrate create -ext sql -dir migrations -seq create_urls
```

Apply migrations:
```bash
migrate -path ./migrations -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable" up
```
