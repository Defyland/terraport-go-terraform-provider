# Senior Readiness Spec

## Product Bar

Terraport must read as a believable Terraform provider for platform teams that manage fictitious BankPort integration resources through infrastructure as code.

## Domain Bar

The domain must use provider-specific language: provider configuration, API product, partner application, webhook endpoint, rate-limit policy, sandbox environment, Terraform state, drift, import, and sensitive generated values.

## Architecture Bar

The repository must show Terraform Plugin Framework boundaries, fake BankPort API boundaries, retry/timeout behavior, and the decision not to use CDKTF as the core implementation.

## API Bar

The fake BankPort API must be documented with OpenAPI, versioned endpoints, auth, error payloads, request/response examples, retryable statuses, and lifecycle endpoint coverage.

## Data and Consistency Bar

The provider must document that Terraform state is a cache and the remote API is the source of truth. Drift detection must be tested by mutating remote state between Terraform refresh/plan steps.

## Security Bar

Token configuration must support provider block and environment variables. Generated secrets must use `Sensitive` schema fields. Docs must explicitly state that sensitive values can still exist in Terraform state and require encrypted remote state.

## Observability Bar

The client must expose request and retry counters for tests and benchmarks. Provider diagnostics must report sanitized status and error codes without logging secrets.

## Performance Bar

The repository must include a repeatable benchmark path for 100 fake resources and retry backoff measurement. Results must record the command used and any local environment caveats.

## Scalability Bar

Docs must identify Terraform refresh/apply remote-call volume, retry amplification, API rate limits, and where batching would become necessary.

## Operational Cost Bar

Docs must discuss the cost of running no service, the hidden cost of Terraform state hygiene, fake API test maintenance, and user support during drift/import incidents.

## Maintainability Bar

Provider code must be split into HTTP client, provider configuration, resources, data sources, helpers, tests, and examples. Adding a resource should follow a clear existing pattern.

## Readability Bar

Code, docs, examples, and tests must use BankPort platform nouns rather than generic Terraform placeholders.

## Test and CI Bar

The repository must include unit/client tests, fake-server acceptance tests, import/drift/failure coverage, benchmarks, formatting/build commands, and CI workflow definitions.

## Evidence Matrix

| Criterion | Evidence | Status | Notes |
| --- | --- | --- | --- |
| Product problem is explicit | README.md, docs/product/problem.md | Done | Names partner onboarding and Terraform workflow. |
| Provider lifecycle is implemented | internal/provider/*resource*.go | Done | CRUD and import implemented for four resources. |
| Fake API client exists | internal/bankport/client.go | Done | Includes retry, timeout, metrics, and sanitized errors. |
| Acceptance tests cover lifecycle and failures | internal/provider/provider_acceptance_test.go | Done | Runs against `httptest.Server` through Terraform Plugin Testing. |
| Sensitive values are protected in schema | internal/provider/*resource*.go, docs/security/secrets.md | Done | Sensitive state leakage is documented as residual risk. |
| Drift detection is verified | internal/provider/provider_acceptance_test.go | Done | Remote mutation produces non-empty Terraform plan. |
| 100-resource performance path exists | internal/provider/provider_acceptance_test.go, benchmarks/baseline.md | Done | `TestAccHundredPartnerAppsApply` and native benchmarks recorded. |
| CDKTF is not core | docs/adr/0002-do-not-use-cdktf-as-core.md | Done | CDKTF retained only as rejected alternative. |
| Runbooks cover expected incidents | docs/runbooks/*.md | Done | Auth, drift, import, timeout, rate limit, state leak. |
| CI validates code and docs | .github/workflows/ci.yml | Done | Includes fmt, test, build, acceptance, security, OpenAPI, benchmark, Docker build. |

## Out of Scope

- Publishing to the Terraform Registry.
- Real BankPort, PixGuard, or SettleFlow production APIs.
- CDKTF constructs or TypeScript generation.
- Long-lived local database or daemon process.
