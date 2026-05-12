# ADR 002: Use PostgreSQL With PostGIS

Status: Accepted

TerraRoute uses PostgreSQL with PostGIS.

## Context

The MVP stores operational data and GPS coordinates. Location queries should use proven geospatial database capabilities.

## Decision

Use PostgreSQL as the primary database and enable PostGIS from the first migration.

## Consequences

- A single relational database supports both business and geospatial data.
- All schema changes must use migrations.
