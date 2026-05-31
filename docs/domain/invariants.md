# Invariants

- Provider `token` is sensitive and can come from the provider block or environment.
- Generated `client_secret`, `signing_secret`, and `api_key_token` are sensitive in schema.
- Terraform state may still contain sensitive values, so state backends must be encrypted and access-controlled.
- `429` and `5xx` responses are retryable; `401` and `404` are not.
- `404` during read removes the Terraform resource from state because the remote API no longer has it.
- Import uses the remote resource ID and then calls read to hydrate state.
- Plan-only resource configuration should not call the remote API before apply.
