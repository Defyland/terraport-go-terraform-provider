# Runbook - API Rate Limited

## Symptom

Terraform diagnostics show `status=429 code=rate_limited`, or apply takes longer than expected due to retries.

## Checks

1. Review provider `retry_max_attempts` and `retry_min_delay_ms`.
2. Reduce Terraform parallelism with `terraform apply -parallelism=5` for large workspaces.
3. Split unrelated partner onboarding changes into smaller applies.
4. Check whether a remote API incident is returning broad `429` responses.

## Recovery

Rerun plan after the rate-limit window clears. Avoid raising retry attempts during a platform-wide incident because retries can amplify load.

## Evidence

`TestAccRateLimitRetry` forces two fake `429` responses and verifies the third create attempt succeeds.
