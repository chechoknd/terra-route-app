# Database

TerraRoute uses PostgreSQL with PostGIS.

## Local Database

Docker Compose starts `postgis/postgis:16-3.4`.

Default local values are defined in `.env.example`.

## Migrations

Migration files live in `backend/migrations`.

Run migrations locally with:

```sh
make migrate-up
```

The initial migration enables:

- `postgis`
- `pgcrypto`

## Current Schema

### companies

Tenant root table for transport companies.

Columns:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `name TEXT NOT NULL`
- `slug TEXT NOT NULL`
- `status TEXT NOT NULL DEFAULT 'active'`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`

Constraints and indexes:

- `status` must be `active`, `inactive`, or `suspended`
- unique case-insensitive `slug`
- index on `status`

### users

Company-scoped user table for future authentication and authorization.

Columns:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `company_id UUID REFERENCES companies(id) ON DELETE RESTRICT`
- `email TEXT NOT NULL`
- `full_name TEXT NOT NULL`
- `role TEXT NOT NULL`
- `status TEXT NOT NULL DEFAULT 'active'`
- `password_hash TEXT NOT NULL`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`

Constraints and indexes:

- `role` must be `super_admin`, `company_admin`, `operator`, or `driver`
- `status` must be `active`, `inactive`, or `suspended`
- `company_id` is required for all roles except `super_admin`
- `super_admin` must not have `company_id`
- unique case-insensitive `email`
- indexes on `company_id`, `(company_id, role)`, and `(company_id, status)`

## Planned MVP Tables

- companies
- users
- vehicles
- drivers
- driver_events
- routes
- route_stops
- trips
- vehicle_locations
- vehicle_last_locations
- incidents

All business tables must include `company_id` unless the table represents companies themselves.
