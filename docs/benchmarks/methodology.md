# Benchmark Methodology

## Goals

- Measure local fake API client cost for 100 partner app creates.
- Measure retry overhead for two `429` responses followed by success.
- Verify Terraform provider quality paths for 100 resources and plan-only remote call avoidance.

## Commands

```sh
make test
make bench
```

## Dataset

- 100 partner applications with one redirect URI and one scope.
- Retry scenario uses two forced `429` responses and one successful product metadata response.

## Metrics Captured

- `ns/op`
- bytes allocated per operation
- allocation count
- provider fake API request counts in tests

## k6 Notes

k6 scripts are included for a BankPort-compatible HTTP API endpoint. They are not the primary benchmark for the provider because Terraport does not expose an HTTP server.
