# TerraRoute

TerraRoute is a SaaS platform for regional and intermunicipal transport companies.

The MVP focuses on secure login, vehicles, drivers, routes, trips, GPS tracking, incidents, and an operational dashboard.

## Monorepo Structure

```text
.
├── backend/          # Go REST API and WebSocket backend
├── admin-web/        # Angular admin dashboard placeholder
├── mobile-driver/    # Flutter driver app placeholder
├── docs/             # Internal project documentation
├── infrastructure/   # Infrastructure notes and future IaC
├── scripts/          # Local development scripts
├── prompts/          # AI-agent prompts
└── docker-compose.yml
```

## Local Setup

1. Copy the environment template:

   ```sh
   cp .env.example .env
   ```

2. Start PostGIS and the API:

   ```sh
   docker compose up --build
   ```

3. Check the API:

   ```sh
   curl http://localhost:18080/healthz
   curl http://localhost:18080/readyz
   ```

4. Apply migrations:

   ```sh
   make migrate-up
   ```

5. Optional: seed local-only demo data:

   ```sh
   make seed-local
   ```

   Demo credentials are documented in `docs/dev-setup.md`.

## Backend

The backend is a Go modular monolith using lightweight clean architecture. It exposes versioned REST endpoints under `/api/v1` and will later include WebSocket location streams.

Current endpoints:

- `GET /healthz`
- `GET /readyz`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`

## Development Rules

- Keep the MVP simple.
- Do not add microservices, Kafka, Kubernetes, CQRS, or event sourcing.
- Every business entity must include `company_id`.
- Every company-scoped query must enforce tenant isolation.
- Never commit secrets. Use `.env.example` for placeholders.

## Documentation

Architecture and database changes must update the relevant files in `docs/`.
