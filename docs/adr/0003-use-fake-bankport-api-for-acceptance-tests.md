# ADR 0003 - Use Fake BankPort API for Acceptance Tests

## Status

Accepted

## Context

The provider needs acceptance tests for create, update, import, delete, drift, auth failure, timeout, and rate-limit retry. No real BankPort control plane exists.

## Options Considered

1. Mock the client and unit test resources only.
2. Use `httptest.Server` as a fake BankPort API.
3. Require a shared external sandbox.

## Decision

Use an in-memory `httptest.Server` fake API in provider tests.

## Consequences

Positive:
- Tests run locally and in CI without credentials.
- The provider still exercises HTTP, auth headers, retries, Terraform lifecycle, and import.

Negative:
- Fake API compatibility does not prove real API compatibility.
- The fake server must evolve with provider schemas.
