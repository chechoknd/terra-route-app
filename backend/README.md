# TerraRoute Backend

Go API for the TerraRoute modular monolith.

## Commands

```sh
go run ./cmd/api
go test ./...
go build ./cmd/api
```

## Current Endpoints

- `GET /healthz`: process health check
- `GET /readyz`: dependency readiness check

## Structure

- `cmd/api`: application entrypoint
- `internal/config`: environment configuration
- `internal/database`: PostgreSQL connection setup
- `internal/server`: HTTP server and routes
- `internal/*`: future MVP modules
- `migrations`: database migration files
