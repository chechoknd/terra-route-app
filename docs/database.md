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
