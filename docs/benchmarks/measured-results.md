# Measured Results

See [../../benchmarks/baseline.md](../../benchmarks/baseline.md) for the current baseline.

## Summary

| Scenario | Evidence | Result |
| --- | --- | --- |
| 100 fake partner app creates | `BenchmarkClientCreate100PartnerApps` | 11.49 ms/op on Apple M1 Max |
| Two rate-limit retries | `BenchmarkClientRetry429Twice` | 4.13 ms/op including backoff |
| 100 Terraform resources | `TestAccHundredPartnerAppsApply` | Passes in fake API acceptance suite |
| Plan-only remote calls | `TestAccPlanOnlyAvoidsRemoteResourceCalls` | Zero fake API requests |

## Next Optimization

If large Terraform workspaces show slow refresh, the next step is a platform API batch-read endpoint. Provider-local caching is intentionally avoided because it can hide drift.
