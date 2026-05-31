# Runbook - Resource Drift Detected

## Symptom

Terraform plan shows changes even though HCL has not changed.

## Checks

1. Identify the resource and attribute changed by the refresh step.
2. Check platform audit logs for manual changes.
3. Decide whether Terraform configuration or remote state is correct.
4. If Terraform is correct, apply to reconcile the remote API.
5. If remote state is correct, update HCL and review the change.

## Recovery

Use `terraform plan` to confirm the intended direction, then apply or update code. Do not edit Terraform state directly unless the remote resource no longer exists.

## Evidence

`TestAccPartnerAppLifecycleImportDrift` mutates the fake API after apply and verifies Terraform returns a non-empty plan.
