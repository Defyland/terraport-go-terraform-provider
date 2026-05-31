# Runbook - Import Existing Resource

## Symptom

A partner app, webhook, rate policy, or sandbox already exists outside Terraform.

## Procedure

1. Find the remote resource ID in the BankPort platform.
2. Add matching HCL configuration.
3. Run import:

```sh
terraform import terraport_bankport_partner_app.main app_0001
```

4. Run `terraform plan`.
5. Resolve any differences between HCL and remote state before applying.

## Notes

Secret version counters are Terraform configuration concerns. Imported resources may need `client_secret_version = 1` or `signing_secret_version = 1` unless a rotation is intentionally planned.
