# Implementation Plan

## Scope

Build a self-contained Terraform provider repository in Go with a fake BankPort control-plane API used by acceptance tests. The first implementation covers provider configuration, one data source, four resources, drift/import behavior, sensitive values, retries, benchmarks, examples, CI, and senior-level documentation.

## Files to Create or Update

- `cmd/terraform-provider-terraport/main.go`
- `internal/bankport/*.go`
- `internal/provider/*.go`
- `internal/provider/*_test.go`
- `examples/**/*.tf`
- `docs/**/*.md`
- `benchmarks/**`
- `openapi.yaml`
- `.github/workflows/ci.yml`
- `Makefile`

## Acceptance Criteria Mapping

| Acceptance criterion | Planned evidence |
| --- | --- |
| Provider config supports endpoint, token, env fallback, timeout, and custom endpoint | `internal/provider/provider.go`, acceptance tests |
| Fake API client for BankPort/PixGuard/SettleFlow-style platform APIs | `internal/bankport/client.go`, fake test server |
| Main resources implement CRUD and import | Resource files and `TestAcc*` tests |
| Data source resolves API product metadata | `api_product_data_source.go`, acceptance test |
| Sensitive generated values are marked and documented | Resource schemas, `docs/security/secrets.md` |
| Drift is detected by Terraform plan | Acceptance test mutating fake API state |
| Retry/backoff handles 429/5xx | Client tests and acceptance rate-limit test |
| Timeout is configurable | Client timeout test and provider config |
| Examples cover required resources | `examples/partner_app`, `examples/webhook_endpoint`, `examples/rate_limit_policy`, `examples/sandbox_environment` |
| 100-resource quality evidence exists | Provider performance test and benchmark docs |
| Docs satisfy general project spec | README, docs tree, OpenAPI, ADRs, runbooks |

## Verification Commands

```sh
make fmt
make test
make test-acceptance
make bench
make build
```

## Risks

- Terraform Plugin Testing requires a local Terraform binary. The repository pins `.tool-versions` for local asdf users and CI installs Terraform explicitly.
- Terraform `Sensitive` masks CLI output but does not remove secrets from state.
- Fake API tests prove provider behavior, not real external API compatibility.

## Deferred Work

- Registry publishing metadata and generated provider website docs.
- Real BankPort API authentication scopes and pagination.
- PixGuard and SettleFlow-specific resources beyond product metadata naming.
