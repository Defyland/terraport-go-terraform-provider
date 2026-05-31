# ADR 0005 - Support Import for Managed Resources

## Status

Accepted

## Context

Platform teams often have resources created manually before Terraform adoption.

## Options Considered

1. Force all resources to be recreated under Terraform.
2. Support import by ID and read remote state.
3. Build a custom migration CLI.

## Decision

Each managed resource implements `ImportState` using the remote resource ID.

## Consequences

Positive:
- Existing partner resources can move under Terraform management.
- Import tests verify state hydration through the fake API.

Negative:
- Import cannot recover secret rotation version counters that exist only in Terraform configuration.
- Imported sensitive values depend on whether the remote API returns them.
