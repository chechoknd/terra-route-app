# TerraRoute Agent Rules

## Project Overview

TerraRoute is a SaaS platform for regional and intermunicipal transport companies.

Main MVP features:
- vehicles
- drivers
- routes
- trips
- GPS tracking
- incidents
- operational dashboard

## Architecture

- Monorepo architecture
- Modular monolith backend
- Go backend
- Angular admin dashboard
- Flutter mobile app
- PostgreSQL + PostGIS
- REST API + WebSockets

## Critical Rules

DO NOT:
- implement microservices
- add Kafka
- add Kubernetes
- overengineer
- add features outside MVP

## Security Rules

# BRANCH WORKFLOW RULES

CRITICAL:

NEVER work directly on `main`.

Before making any code changes, always:

1. Check the current branch.
2. If the current branch is `main`, stop and create a feature branch.
3. Use a branch name that follows the project convention.

Branch naming:

- feature/module-name
- fix/module-name
- refactor/module-name
- chore/module-name

Examples:

- chore/sprint-0-foundation
- feature/auth-companies-users
- feature/vehicles-crud
- fix/docker-postgis-config

All implementation work must happen in a non-main branch.

The `main` branch must only receive changes through Pull Requests.

NEVER commit:
- `.env` files
- secrets
- API keys
- credentials
- tokens

Always use:
- `.env.example`

## Multi-Tenancy Rules

Every business entity must contain:
- `company_id`

Every company-scoped query must validate tenant scope. Never expose data between companies.

This is CRITICAL.

## Backend Rules

- Follow idiomatic Go.
- Keep modules isolated.
- Prefer readability.
- Use migrations for schema changes.
- Add tests for business logic.
- Use `context.Context` for request-scoped work.

## Frontend Rules

- Use Angular standalone components.
- Use NgRx.
- Use feature-based structure.

## Mobile Rules

- Flutter only.
- GPS background tracking.
- Offline-first mindset.

## Documentation Rules

Update docs when architecture changes.

Main docs:
- `docs/architecture.md`
- `docs/database.md`
- `docs/api.md`
- `docs/project-status.md`

## Git Rules

Branch naming:
- `feature/module-name`
- `fix/module-name`
- `refactor/module-name`

Commit format:
- `feat(module): description`
- `fix(module): description`
