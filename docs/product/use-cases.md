# Use Cases

## Create a Partner Application

Define redirect URIs and scopes in Terraform. The provider returns `client_id` and a sensitive `client_secret`.

## Attach a Webhook Endpoint

Create a webhook endpoint tied to a partner app, receive a sensitive `signing_secret`, and rotate it by incrementing `signing_secret_version`.

## Govern API Consumption

Attach a rate-limit policy to a partner app or product subject before the partner starts load testing.

## Provision Sandbox Access

Provision a sandbox environment for BankPort, PixGuard, and SettleFlow-style products and return a sensitive `api_key_token`.

## Import Existing Platform Resources

Bring manually created resources into Terraform state using resource IDs and immediately detect whether configuration matches remote state.
