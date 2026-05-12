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

### vehicles

Company-scoped vehicle table for fleet inventory.

Columns:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `company_id UUID NOT NULL REFERENCES companies(id) ON DELETE RESTRICT`
- `plate TEXT NOT NULL`
- `internal_code TEXT NOT NULL`
- `vehicle_type TEXT NOT NULL`
- `brand TEXT NOT NULL`
- `model TEXT NOT NULL`
- `capacity INTEGER NOT NULL`
- `status TEXT NOT NULL DEFAULT 'active'`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`

Constraints and indexes:

- `company_id` is required for tenant isolation
- `capacity` must be greater than `0`
- `status` must be `active`, `inactive`, `maintenance`, or `unavailable`
- unique case-insensitive `plate` per company
- unique case-insensitive `internal_code` per company
- indexes on `company_id`, `(company_id, status)`, and `(company_id, vehicle_type)`

### drivers

Company-scoped driver table for operational driver records.

Columns:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `company_id UUID NOT NULL REFERENCES companies(id) ON DELETE RESTRICT`
- `user_id UUID REFERENCES users(id) ON DELETE SET NULL`
- `first_name TEXT NOT NULL`
- `last_name TEXT NOT NULL`
- `document_number TEXT NOT NULL`
- `phone TEXT NOT NULL`
- `email TEXT`
- `license_number TEXT NOT NULL`
- `status TEXT NOT NULL DEFAULT 'active'`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`

Constraints and indexes:

- `company_id` is required for tenant isolation
- `user_id` is optional and links a driver record to an application user when needed
- `status` must be `active`, `inactive`, or `suspended`
- unique case-insensitive `document_number` per company
- unique case-insensitive `email` per company when email is present
- indexes on `company_id`, `(company_id, status)`, `(company_id, user_id)`, and `(company_id, lower(license_number))`

### routes

Company-scoped route table for intermunicipal route definitions.

Columns:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `company_id UUID NOT NULL REFERENCES companies(id) ON DELETE RESTRICT`
- `name TEXT NOT NULL`
- `origin_city TEXT NOT NULL`
- `destination_city TEXT NOT NULL`
- `estimated_distance_km NUMERIC(10,2) NOT NULL`
- `estimated_duration_minutes INTEGER NOT NULL`
- `base_price NUMERIC(12,2) NOT NULL`
- `status TEXT NOT NULL DEFAULT 'active'`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`

Constraints and indexes:

- `company_id` is required for tenant isolation
- `name`, `origin_city`, and `destination_city` must not be empty
- `estimated_distance_km` must be `0` or greater
- `estimated_duration_minutes` must be greater than `0`
- `base_price` must be `0` or greater
- `status` must be `active`, `inactive`, or `archived`
- unique case-insensitive `name` per company
- indexes on `company_id`, `(company_id, status)`, and `(company_id, lower(origin_city), lower(destination_city))`

### route_stops

Ordered stops that belong to a route.

Columns:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE`
- `name TEXT NOT NULL`
- `city TEXT NOT NULL`
- `stop_order INTEGER NOT NULL`
- `latitude DOUBLE PRECISION NOT NULL`
- `longitude DOUBLE PRECISION NOT NULL`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`

Constraints and indexes:

- route stops are scoped through their parent route; repository queries join `routes` and validate `routes.company_id`
- `name` and `city` must not be empty
- `stop_order` must be greater than `0`
- `latitude` must be between `-90` and `90`
- `longitude` must be between `-180` and `180`
- unique `stop_order` per route
- indexes on `route_id` and `(route_id, lower(city))`

## Planned MVP Tables

- companies
- users
- vehicles
- drivers
- driver_events
- trips
- vehicle_locations
- vehicle_last_locations
- incidents

All business tables must include `company_id` unless the table represents companies themselves.
