# Deployment View

Terraport is distributed as a local Terraform provider binary. It does not deploy a service.

## Local Development

1. Build with `make build`.
2. Run tests with fake API using `make test`.
3. Use examples with a BankPort-compatible endpoint.

## CI

GitHub Actions installs Go and Terraform, runs formatting checks, builds the provider, runs fake API tests, benchmarks, OpenAPI lint, security scanning, and Docker build validation.

## Public Registry Path

Publishing is intentionally deferred. A real publishing workflow would add signing, release notes, provider documentation generation, and Terraform Registry metadata.
