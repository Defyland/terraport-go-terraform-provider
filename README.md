# Terraport Terraform Provider

Terraport is a Terraform provider written in Go for provisioning fictitious platform resources used by the backend challenge portfolio: BankPort partner applications, webhook endpoints, rate-limit policies, and sandbox environments. It is intentionally built as a provider project, not a CDKTF application, so the engineering evidence focuses on Terraform Plugin Framework lifecycle behavior, remote platform APIs, drift, import, sensitive state, retries, and acceptance tests.

## 1. What is this product?

Terraport lets platform engineers declare BankPort integration infrastructure as code. A partner team can create an OAuth-style application, attach webhooks, define rate-limit controls, and provision a sandbox environment from Terraform.

## 2. Problem it solves

Partner onboarding normally mixes manual portal clicks, copied secrets, and undocumented rate-limit exceptions. Terraport moves those platform resources into reviewed Terraform changes with import support and drift detection.

## 3. Target users

- Platform engineers who operate shared BankPort, PixGuard, and SettleFlow sandbox APIs.
- Partner integration teams that need repeatable onboarding environments.
- Security reviewers who need evidence that generated secrets are marked sensitive and documented.

## 4. Main features

- Provider configuration through `endpoint`, `token`, environment variables, retry settings, and request timeouts.
- Fake BankPort HTTP API client for deterministic acceptance tests with `httptest`.
- Resources: `terraport_bankport_partner_app`, `terraport_bankport_webhook_endpoint`, `terraport_bankport_rate_limit_policy`, and `terraport_bankport_sandbox_environment`.
- Data source: `terraport_bankport_api_product`.
- Full create, read, update, delete, and import lifecycle for resources.
- Retry and exponential backoff for `429` and `5xx` responses.
- Sensitive schemas for `client_secret`, `signing_secret`, and `api_key_token`.

## 5. Architecture overview

Terraform Core talks to the provider through Terraform Plugin Framework protocol v6. The provider translates resource lifecycle calls into HTTP requests against the BankPort platform API. Tests replace the real platform with an in-memory `httptest.Server`, which gives repeatable acceptance tests without mocking Terraform Core.

See [docs/architecture/overview.md](docs/architecture/overview.md) and [docs/engineering-case-study.md](docs/engineering-case-study.md).

## 6. Tech stack

- Go 1.25
- Terraform Plugin Framework
- Terraform Plugin Testing
- `httptest` fake platform API
- GitHub Actions for build, tests, formatting, and security workflow checks

## 7. Domain model

The provider models platform API products, partner applications, webhook endpoints, rate-limit policies, and sandbox environments. The remote BankPort API is the source of truth; Terraform state is a cached desired/observed representation used for plans.

See [docs/domain/glossary.md](docs/domain/glossary.md).

## 8. API documentation

The provider consumes the documented fake BankPort API in [openapi.yaml](openapi.yaml). Terraform users normally do not call this API directly, but the contract is documented so tests and provider behavior are auditable.

## 9. Async or event architecture

The provider does not run background workers or publish broker events. The fake API documents platform audit events in [docs/events/README.md](docs/events/README.md) because resource lifecycle operations would create audit events in a real BankPort control plane.

## 10. Database design

The provider has no local database. Terraform state stores resource identity, non-secret attributes, and sensitive computed values marked with Terraform's `Sensitive` schema flag. The fake API uses an in-memory store only during tests.

## 11. Testing strategy

- Client tests cover retry, timeout, authentication failure, and redaction behavior.
- Acceptance tests run Terraform against an `httptest` BankPort API for create, update, import, delete, drift, auth failure, timeout, and rate-limit retry cases.
- Benchmarks measure the fake client and a 100-resource apply path.

## 12. Performance benchmarks

Benchmark methodology and current measured output live in [docs/benchmarks/methodology.md](docs/benchmarks/methodology.md) and [benchmarks/baseline.md](benchmarks/baseline.md). Results are intentionally modest: the point is to show remote call counts, retry cost, and state refresh behavior.

## 13. Observability

Provider diagnostics include sanitized API status, retry metrics are exposed through the client for tests and benchmarks, and the fake API records request counts by method/path. The provider avoids logging token and generated secrets.

## 14. Security considerations

Secrets are marked `Sensitive`, but Terraform still stores them in state. This is documented in [docs/security/secrets.md](docs/security/secrets.md) and [docs/runbooks/state-sensitive-leak-risk.md](docs/runbooks/state-sensitive-leak-risk.md).

## 15. Trade-offs and decisions

ADRs cover Terraform Plugin Framework, rejecting CDKTF as core implementation, fake API acceptance tests, sensitive generated values, import support, and retry/timeouts. See [docs/adr](docs/adr).

## 16. How to run locally

```sh
make build
```

For manual Terraform experiments, build the provider binary and use the examples under [examples](examples). The examples expect a BankPort-compatible endpoint and token.

## 17. How to run tests

```sh
make test
make test-acceptance
make bench
```

The acceptance tests use a local fake server and do not require real BankPort credentials.

## 18. Failure scenarios

Documented runbooks cover provider authentication failure, remote drift, importing existing resources, apply timeout, API rate limiting, and Terraform state secret exposure risk.

## 19. Roadmap

- Add PixGuard and SettleFlow product-specific resources once their platform control-plane APIs are stable.
- Add plan modifiers for more precise replacement behavior on immutable remote fields.
- Add structured debug logging with explicit redaction hooks.
- Publish provider docs generated from schema when the public registry layout is introduced.
