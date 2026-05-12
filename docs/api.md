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

## Authorization

Role authorization is enforced from validated JWT claims stored in request context. Roles are never accepted from request bodies.

Supported roles:

- `super_admin`
- `company_admin`
- `operator`
- `driver`

Authorization errors use JSON responses:

```json
{
  "error": "forbidden"
}
```

### Vehicles

Vehicle endpoints require a valid `Authorization: Bearer <access_token>` header and are limited to `company_admin` and `operator`.

Drivers cannot manage vehicles.

All vehicle operations use `company_id` from the authenticated JWT claims. `company_id` from request bodies is ignored.

```http
GET /api/v1/vehicles
```

Returns vehicles for the authenticated company.

```http
POST /api/v1/vehicles
```

Creates a vehicle for the authenticated company.

Request body:

```json
{
  "plate": "ABC123",
  "internal_code": "BUS-001",
  "vehicle_type": "bus",
  "brand": "Mercedes-Benz",
  "model": "OF-1721",
  "capacity": 42,
  "status": "active"
}
```

```http
GET /api/v1/vehicles/{id}
PATCH /api/v1/vehicles/{id}
DELETE /api/v1/vehicles/{id}
```

`DELETE` marks the vehicle inactive instead of hard deleting it.

Vehicle errors use JSON responses:

```json
{
  "error": "vehicle_not_found"
}
```

### Drivers

Driver endpoints require a valid `Authorization: Bearer <access_token>` header and are limited to `company_admin` and `operator`.

Drivers cannot manage driver records.

All driver operations use `company_id` from the authenticated JWT claims. `company_id` from request bodies is ignored.

```http
GET /api/v1/drivers
```

Returns drivers for the authenticated company.

```http
POST /api/v1/drivers
```

Creates a driver for the authenticated company.

Request body:

```json
{
  "user_id": "optional linked user uuid",
  "first_name": "Ana",
  "last_name": "Torres",
  "document_number": "DOC-001",
  "phone": "+573001112233",
  "email": "ana@example.test",
  "license_number": "LIC-001",
  "status": "active"
}
```

```http
GET /api/v1/drivers/{id}
PATCH /api/v1/drivers/{id}
DELETE /api/v1/drivers/{id}
```

`DELETE` marks the driver inactive instead of hard deleting it.

Driver responses do not expose password hashes or user credential fields.

Driver errors use JSON responses:

```json
{
  "error": "driver_not_found"
}
```

### Routes

Route endpoints require a valid `Authorization: Bearer <access_token>` header and are limited to `company_admin` and `operator`.

All route operations use `company_id` from the authenticated JWT claims. `company_id` from request bodies is ignored.

Route stops are not managed by these endpoints yet.

```http
GET /api/v1/routes
```

Returns routes for the authenticated company.

```http
POST /api/v1/routes
```

Creates a route for the authenticated company.

Request body:

```json
{
  "name": "Bogota - Tunja",
  "origin_city": "Bogota",
  "destination_city": "Tunja",
  "estimated_distance_km": 140.5,
  "estimated_duration_minutes": 180,
  "base_price": 45000,
  "status": "active"
}
```

```http
GET /api/v1/routes/{id}
PATCH /api/v1/routes/{id}
DELETE /api/v1/routes/{id}
```

`DELETE` archives the route instead of hard deleting it.

Route errors use JSON responses:

```json
{
  "error": "route_not_found"
}
```

### Route Stops

Route stop endpoints require a valid `Authorization: Bearer <access_token>` header and are limited to `company_admin` and `operator`.

Access is scoped through the parent route. The route must belong to the authenticated company.
If the parent route does not belong to the authenticated company, the API returns `route_stop_not_found`.

```http
GET /api/v1/routes/{id}/stops
```

Returns stops ordered by `stop_order`.

```http
POST /api/v1/routes/{id}/stops
```

Adds a stop to the route.

Request body:

```json
{
  "name": "Terminal Norte",
  "city": "Bogota",
  "stop_order": 1,
  "latitude": 4.710989,
  "longitude": -74.072092
}
```

```http
PATCH /api/v1/routes/{id}/stops/{stopId}
DELETE /api/v1/routes/{id}/stops/{stopId}
```

Route stop errors use JSON responses:

```json
{
  "error": "route_stop_not_found"
}
```
