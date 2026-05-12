# Development Setup

## Local Seed Data

TerraRoute includes local-only seed data for development. This seed is idempotent and can be run multiple times after migrations.

Run:

```sh
make migrate-up
make seed-local
```

The seed creates:

- one demo company: `TerraRoute Demo Company`
- one `company_admin` user
- one `operator` user
- one `driver` user

Local-only demo credentials:

| Role | Email | Password |
| --- | --- | --- |
| company_admin | `admin@terraroute.local` | `TerraRoute123!` |
| operator | `operator@terraroute.local` | `TerraRoute123!` |
| driver | `driver@terraroute.local` | `TerraRoute123!` |

These are not real credentials. Do not use them outside local development.

The seed is not part of migrations and is not run automatically by Docker Compose or the API.
