# ADR 003: Use Monorepo Architecture

Status: Accepted

TerraRoute uses a monorepo for backend, web, mobile, docs, prompts, and infrastructure.

## Context

The project is built incrementally with AI coding agents. A monorepo makes project context easier to discover and reduces coordination overhead.

## Decision

Keep all MVP applications and docs in one repository.

## Consequences

- Shared documentation and prompts are close to code.
- Agents must avoid unrelated edits across app boundaries.
