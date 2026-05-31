# ADR 0001 - Use Terraform Plugin Framework

## Status

Accepted

## Context

The project must demonstrate provider lifecycle behavior, state, import, drift, sensitive attributes, and Terraform-native testing.

## Options Considered

1. Terraform Plugin SDK v2.
2. Terraform Plugin Framework.
3. A custom CLI that writes Terraform JSON.

## Decision

Use Terraform Plugin Framework and protocol v6.

## Consequences

Positive:
- Resource and data source implementations use Terraform's current strongly typed framework.
- Sensitive attributes and import behavior are modeled in provider schemas.
- Terraform Plugin Testing can run fake API lifecycle tests.

Negative:
- Framework models require explicit conversion between Terraform values and Go structs.
- Acceptance tests need Terraform-compatible test harness setup.
