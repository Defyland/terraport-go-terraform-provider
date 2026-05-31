# ADR 0004 - Mark Generated Secrets Sensitive

## Status

Accepted

## Context

Partner apps, webhooks, and sandboxes return generated secrets. Terraform can mask sensitive output, but state may still contain the raw value.

## Options Considered

1. Do not return generated secrets.
2. Return generated secrets as normal attributes.
3. Return generated secrets as `Sensitive` computed attributes and document state risk.

## Decision

Use sensitive computed attributes for `client_secret`, `signing_secret`, and `api_key_token`.

## Consequences

Positive:
- Terraform CLI output masks the values.
- Users can pass the values to other Terraform resources without plain output.

Negative:
- Terraform state still contains the values.
- Users must use encrypted state backends and strict access controls.
