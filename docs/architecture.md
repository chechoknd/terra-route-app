# Architecture

TerraRoute uses a monorepo with three application areas:

- `backend`: Go modular monolith API.
- `admin-web`: Angular admin dashboard.
- `mobile-driver`: Flutter driver app.

The backend follows a lightweight clean architecture inside a modular monolith. Each business module owns its domain, application, infrastructure, and HTTP interface code.

The MVP explicitly avoids microservices, Kafka, Kubernetes, CQRS, event sourcing, and GraphQL.

## Backend Modules

Planned MVP modules:

- auth
- companies
- users
- vehicles
- drivers
- driver_events
- routes
- route_stops
- trips
- locations
- incidents
- dashboard

## Multi-Tenancy

Every business entity must include `company_id`. All company-scoped queries must validate tenant scope.

## Runtime

Local runtime uses Docker Compose with:

- API container
- PostgreSQL with PostGIS

Redis is not part of the MVP foundation.
