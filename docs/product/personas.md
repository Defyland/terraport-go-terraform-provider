# Personas

## Platform Engineer

Owns shared BankPort control-plane resources. Needs reviewed changes, predictable import behavior, and failure diagnostics that do not leak tokens.

## Partner Integration Engineer

Needs sandbox credentials, OAuth client IDs, webhook secrets, and stable rate limits for automated test environments.

## Security Reviewer

Checks scopes, token handling, generated secret state exposure, and whether Terraform state is protected by encryption and access controls.

## Support Engineer

Uses runbooks to diagnose authentication failure, API rate limits, drift, imports, and timeout failures during partner onboarding.
