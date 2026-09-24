# URLShortener


## Table of Contents

- [Description](#description)
- [Configuration](#configuration)
- [First Run](#first-run)


## Description

A URL shortening service written in Go. Uses PostgreSQL for storing URLs.

The service allows users to create short URLs and redirect users using an alias.


## Configuration

Create a `.env` file based on `.env.example`:
```bash
copy .env.example .env
```


## First Run

Start PostgreSQL:
```bash
make up
```

Apply migrations:
```bash
make migrate-up
```

Start the application:
```bash
make run
```

To see all available commands:
```bash
make help
```
