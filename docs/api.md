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

### Authentication

```http
POST /api/v1/auth/login
```

Authenticates a company-scoped user and returns an access token.

Request body:

```json
{
  "company_id": "company uuid",
  "email": "operator@example.com",
  "password": "user password"
}
```

Successful response:

```json
{
  "access_token": "jwt access token",
  "token_type": "Bearer",
  "user": {
    "id": "user uuid",
    "company_id": "company uuid",
    "email": "operator@example.com",
    "full_name": "Operator Name",
    "role": "operator",
    "status": "active",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

```http
GET /api/v1/auth/me
Authorization: Bearer <access_token>
```

Returns the authenticated active user. This endpoint is protected by JWT authentication middleware.

Authentication errors use JSON responses:

```json
{
  "error": "invalid_token"
}
```

## Planned Conventions

```http
GET    /api/v1/vehicles
POST   /api/v1/vehicles
GET    /api/v1/vehicles/{id}
PATCH  /api/v1/vehicles/{id}
DELETE /api/v1/vehicles/{id}
```

Protected endpoints must validate JWT, role, and company scope.
