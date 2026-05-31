# Threat Model

## Assets

- Provider token.
- Generated `client_secret`, `signing_secret`, and `api_key_token`.
- Terraform state files and remote state backend.
- Remote BankPort platform resources.

## Actors

- Authorized platform engineer.
- Read-only reviewer.
- Compromised CI job.
- Unauthorized user with state backend access.
- Operator making manual remote changes outside Terraform.

## Trust Boundaries

- Terraform CLI to provider process.
- Provider process to BankPort API over HTTPS.
- Terraform state backend to humans and CI.
- Test fake API to provider acceptance tests.

## Abuse Cases

- Token copied into HCL instead of environment or secret store.
- Sensitive Terraform state downloaded by unauthorized user.
- Webhook signing secret printed through a non-sensitive output.
- Remote resource changed manually and not detected before deploy.
- API rate limit causes partial apply and repeated retries.

## Controls

- Provider `token` schema is sensitive.
- Generated secrets are sensitive computed attributes.
- Diagnostics pass through redaction before Terraform receives API errors.
- Acceptance tests cover auth failure, drift, timeout, rate limit, import, and lifecycle.
- Runbooks document state leak and drift recovery.

## Residual Risks

- Terraform state contains sensitive values even when schema marks them sensitive.
- The fake API cannot prove real BankPort authorization semantics.
- No provider-side tenant authorization exists; the remote API must enforce it.
