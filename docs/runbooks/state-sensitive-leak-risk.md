# Runbook - State Sensitive Leak Risk

## Symptom

Terraform state may have been exposed, downloaded by an unauthorized identity, or committed by mistake.

## Immediate Actions

1. Treat `client_secret`, `signing_secret`, `api_key_token`, and provider tokens as compromised.
2. Revoke or rotate affected secrets in the BankPort platform.
3. Remove exposed state from public or shared locations.
4. Audit state backend access logs.
5. Re-run Terraform plan after rotation.

## Prevention

- Use encrypted remote state.
- Restrict state backend access to Terraform operators.
- Mark module outputs containing generated secrets as `sensitive = true`.
- Do not commit `.tfstate`, `.tfvars`, or crash logs.

## Important Limitation

Terraform `Sensitive` masks CLI output but does not remove raw values from state.
