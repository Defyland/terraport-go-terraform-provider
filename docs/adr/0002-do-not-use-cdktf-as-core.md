# ADR 0002 - Do Not Use CDKTF as Core

## Status

Accepted

## Context

CDKTF was mentioned historically as an optional direction, but the product goal is a Terraform provider in Go.

## Options Considered

1. Implement a Go Terraform provider.
2. Generate CDKTF constructs and treat Terraform as generated output.
3. Build both provider and CDKTF layer now.

## Decision

Do not use CDKTF as the core implementation. Keep the project focused on Terraform Plugin Framework provider behavior.

## Consequences

Positive:
- The repository demonstrates provider lifecycle and state management directly.
- Drift, import, and sensitive state behavior are tested at the provider layer.

Negative:
- TypeScript users do not get CDKTF constructs in this phase.
- A future CDKTF layer would need to consume provider schemas later.
