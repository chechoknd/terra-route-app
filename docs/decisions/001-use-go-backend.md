# ADR 001: Use Go Backend

Status: Accepted

TerraRoute uses Go for the backend API.

## Context

The MVP needs a simple, reliable backend with good concurrency support for REST APIs and future WebSocket location updates.

## Decision

Use Go with a modular monolith architecture.

## Consequences

- Fast builds and simple deployment.
- Strong standard library support.
- Agents must keep module boundaries explicit and avoid unnecessary abstractions.
