# Benchmark Baseline

## Environment

- Date: 2026-05-31
- Machine: Apple M1 Max
- Go: 1.25.10
- Terraform: 1.9.8 for provider acceptance tests

## Commands

```sh
make test
make bench
```

## Measured Results

```text
BenchmarkClientCreate100PartnerApps-10      100   11489944 ns/op   1077758 B/op   14019 allocs/op
BenchmarkClientRetry429Twice-10             314    4125915 ns/op     28598 B/op     330 allocs/op
```

## 100 Resource Apply Evidence

`TestAccHundredPartnerAppsApply` applies 100 `terraport_bankport_partner_app` resources against the fake API and verifies 100 remote resources exist before test cleanup.

## Plan Remote Call Evidence

`TestAccPlanOnlyAvoidsRemoteResourceCalls` runs a plan-only resource configuration and verifies the fake API receives zero requests. Resource planning should not call the remote API before apply; refresh and data sources are the intentional remote-call paths.

## Retry Evidence

`TestAccRateLimitRetry` forces two `429` responses and verifies the provider succeeds on the third create attempt.

## Bottleneck Found

The first bottleneck is remote API round trips during refresh/apply, not local CPU. Retry settings can multiply remote calls under `429` or `5xx` incidents.
