# Engineering Case Study

## 1. Product Context

Terraport turns BankPort partner onboarding resources into Terraform-managed infrastructure. The product exists for platform teams that want reviewable changes for applications, webhooks, rate limits, and sandbox environments.

## 2. Domain Model

The domain centers on partner apps, webhook endpoints, rate-limit policies, sandbox environments, API products, Terraform state, drift, import, and generated secrets.

## 3. Architecture

Terraform Core loads a Go provider built with Terraform Plugin Framework. Resource lifecycle methods call a BankPort-compatible HTTP client. Acceptance tests replace the real platform with `httptest.Server`.

## 4. Key Trade-offs

The provider uses a fake API instead of a real sandbox because no external product exists. This gives deterministic lifecycle and failure tests, but it does not prove compatibility with a future real BankPort API.

CDKTF is not used as core because the learning target is provider lifecycle, not generated Terraform code.

## 5. Data Model

The provider stores no database. Terraform state contains remote IDs, desired/observed attributes, and sensitive computed values. The remote API is the source of truth for refresh and drift.

## 6. Consistency Model

Operations are request/response HTTP calls. Terraform refresh reads the remote API before comparing configuration. Deletes treat remote `404` as already gone. Import hydrates state by ID through read.

## 7. Failure Scenarios

Covered failures:

- Authentication failure: `TestAccAuthFailure`.
- Timeout: `TestAccApplyTimeout`.
- Rate limit retry: `TestAccRateLimitRetry`.
- Drift: `TestAccPartnerAppLifecycleImportDrift`.
- Missing remote resource: read removes state.

## 8. Performance Strategy

The provider avoids remote calls during plan-only creation. Refresh and apply call the remote API intentionally. Benchmarks measure 100 fake creates and retry overhead.

## 9. Scalability Strategy

The primary scaling limit is remote API calls during refresh/apply. Large workspaces should lower Terraform parallelism during rate-limit incidents. Future real APIs may need batch reads.

## 10. Security Model

Provider token is sensitive. Generated secrets are sensitive in schema. Diagnostics redact secret-like values. The major residual risk is Terraform state exposure.

## 11. Observability

The client exposes metrics counters used by tests and benchmarks: requests, retries, rate-limit responses, and server-error responses. Fake API tests count method/path calls for plan and retry assertions.

## 12. Operational Cost

No service is deployed, but there is still maintenance cost: schema compatibility, test fake API upkeep, Terraform state hygiene, and runbook support for import/drift incidents.

## 13. Maintainability

Code boundaries keep Terraform concerns in `internal/provider` and HTTP concerns in `internal/bankport`. Adding a resource requires a type, client methods, resource implementation, fake API routes, tests, docs, and examples.

## 14. Product Decisions

The provider starts with BankPort resources because partner onboarding has clear lifecycle and sensitive-value behavior. PixGuard and SettleFlow are represented as API products and sandbox product selections until resource APIs are defined.

## 15. What I Would Do Next

- Add schema documentation generation.
- Add provider upgrade tests for state migrations.
- Add correlation ID headers once the platform API specifies tracing.
- Add real API contract tests when a BankPort sandbox exists.
