# Consistency

## Transaction Boundaries

Each Terraform lifecycle call maps to one remote API operation, except secret rotation, which updates the resource first and then calls a rotation endpoint when the version counter increases.

## Source of Truth

The remote BankPort API is the source of truth. Terraform state is refreshed from the remote API before detecting drift.

## Rollback Strategy

Terraform does not provide automatic remote rollback if a later resource fails. Operators should run `terraform plan` after a partial failure, inspect remote state, and apply again or import completed resources.

## Idempotency Assumption

The fake API does not implement explicit idempotency keys. A real API should support idempotency for create operations to reduce duplicate resource risk after network timeouts.
