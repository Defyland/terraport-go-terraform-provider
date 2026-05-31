# Runbook - Apply Timeout

## Symptom

Terraform fails with a timeout or deadline exceeded error during create, read, update, or delete.

## Checks

1. Confirm `timeout_ms` in the provider block.
2. Check BankPort platform latency or incident status.
3. Inspect whether the operation eventually completed remotely.
4. Run `terraform plan` to refresh state before retrying apply.

## Recovery

Increase `timeout_ms` for slow sandboxes or retry after the remote API incident clears. If the remote operation completed but Terraform missed the response, import or refresh the resource before applying again.

## Evidence

`TestAccApplyTimeout` delays the fake API beyond provider timeout and expects Terraform failure.
