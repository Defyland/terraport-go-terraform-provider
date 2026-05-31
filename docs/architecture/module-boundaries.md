# Module Boundaries

## `internal/bankport`

Owns HTTP request construction, response decoding, retry/backoff, timeout behavior, metrics counters, API errors, and JSON types.

It does not import Terraform packages.

## `internal/provider`

Owns Terraform Plugin Framework schemas, provider configuration, data source read, resource lifecycle methods, import behavior, and Terraform diagnostics.

It does not hand-roll HTTP calls.

## `cmd/terraform-provider-terraport`

Owns provider process startup and the provider registry address used by Terraform.

## Test Boundary

The fake API lives in provider tests because it exists to validate provider lifecycle behavior, not as a reusable production server.
