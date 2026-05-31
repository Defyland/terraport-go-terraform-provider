# Bounded Contexts

## Terraform Provider Context

Owns schema, configuration, CRUD lifecycle, import, diagnostics, and state mapping. It does not own remote persistence.

## BankPort Platform API Context

Owns partner app, webhook, rate-limit, sandbox, and product metadata persistence. In this repository the context is represented by `httptest` fake API handlers.

## Terraform State Context

Owns local or remote state storage. The provider only marks attributes sensitive; state encryption and access control belong to the Terraform backend.

## Security Review Context

Owns token issuance, scopes, state access, and incident response for leaked state files.
