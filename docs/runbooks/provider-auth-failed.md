# Runbook - Provider Auth Failed

## Symptom

Terraform apply or refresh fails with `status=401 code=unauthorized`.

## Checks

1. Confirm the provider block does not hard-code an expired token.
2. Check `TERRAPORT_TOKEN` or `BANKPORT_TOKEN` in the shell or CI secret context.
3. Confirm the token has product read and resource write scopes required by the Terraform plan.
4. Re-run `terraform plan` after token rotation.

## Recovery

Rotate the provider token in the platform identity system, update the CI secret or local environment, and rerun plan before apply.

## Prevention

Use short-lived scoped tokens and separate local, CI, and break-glass identities.
