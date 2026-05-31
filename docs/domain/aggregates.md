# Aggregates

## Partner Application

Root: `bankport_partner_app`

Fields: name, product code, redirect URIs, scopes, status, client ID, client secret, client secret version.

Invariant: a partner app must have at least one redirect URI and at least one scope in the Terraform configuration used by examples and tests.

## Webhook Endpoint

Root: `bankport_webhook_endpoint`

Fields: partner app ID, URL, event types, enabled flag, signing secret, signing secret version.

Invariant: webhook endpoint lifecycle depends on a valid partner app ID.

## Rate-Limit Policy

Root: `bankport_rate_limit_policy`

Fields: product code, subject type, subject ID, requests per minute, burst limit, mode.

Invariant: policy mode is operationally either `enforce` or `report`; tests exercise both values.

## Sandbox Environment

Root: `bankport_sandbox_environment`

Fields: name, products, region, status, API key token.

Invariant: the remote API owns token generation and environment ID assignment.
