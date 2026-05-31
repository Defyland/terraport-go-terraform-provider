# Secrets

## Provider Token

The provider supports `token` in the provider block, `TERRAPORT_TOKEN`, and `BANKPORT_TOKEN`. Environment variables are preferred for local use because they avoid committing secrets to HCL.

## Generated Secrets

The provider marks these attributes `Sensitive`:

- `terraport_bankport_partner_app.client_secret`
- `terraport_bankport_webhook_endpoint.signing_secret`
- `terraport_bankport_sandbox_environment.api_key_token`

Terraform masks sensitive values in CLI output, but state can still contain raw values. Use encrypted remote state, strict IAM, short retention, and audit access.

## Rotation

- Partner app secret: increment `client_secret_version`.
- Webhook signing secret: increment `signing_secret_version`.
- Sandbox API key token rotation is deferred because the fake API only generates it during create.
