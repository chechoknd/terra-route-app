# ADR 004: Use Flutter For Driver App

Status: Accepted

TerraRoute uses Flutter for the mobile driver application.

## Context

The driver app must support trip operations, GPS tracking, and offline-friendly behavior.

## Decision

Use Flutter for the mobile app.

## Consequences

- Mobile work remains isolated under `mobile-driver`.
- Native GPS/background behavior must be validated during implementation.
