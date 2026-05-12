# API

TerraRoute uses REST endpoints versioned under `/api/v1`.

## Current Endpoints

### Health

```http
GET /healthz
GET /api/v1/healthz
```

Returns process health.

### Readiness

```http
GET /readyz
```

Checks PostgreSQL connectivity.

## Planned Conventions

```http
GET    /api/v1/vehicles
POST   /api/v1/vehicles
GET    /api/v1/vehicles/{id}
PATCH  /api/v1/vehicles/{id}
DELETE /api/v1/vehicles/{id}
```

Protected endpoints must validate JWT, role, and company scope.
