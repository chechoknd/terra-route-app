# Project Status

Status: Sprint 1 started.

Completed:

- monorepo skeleton
- backend skeleton
- Docker Compose for API and PostGIS
- Dockerized migration runner
- Go HTTP server
- health and readiness endpoints
- structured logger
- graceful shutdown
- migration directory
- initial PostGIS migration
- companies and users schema migration
- company and user domain entities
- company and user PostgreSQL repositories
- bcrypt password hashing service
- login application use case with JWT token response
- JWT generation and validation service
- auth login and me HTTP handlers
- JWT authentication middleware for protected routes
- verified local health/readiness endpoints on port 18080
- internal documentation structure

Next:

- add role authorization middleware
- add companies/users application use cases when backend code needs them
