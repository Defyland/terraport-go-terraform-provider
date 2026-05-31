# API Documentation

The provider consumes the fake BankPort Platform API documented in [../../openapi.yaml](../../openapi.yaml). Tests start an in-memory API with the same routes.

## Authentication

Every request uses:

```http
Authorization: Bearer <token>
Accept: application/json
```

Provider token resolution order:

1. `token` in provider block.
2. `TERRAPORT_TOKEN`.
3. `BANKPORT_TOKEN`.

## Error Format

```json
{
  "error": {
    "code": "rate_limited",
    "message": "fake API rate limit"
  }
}
```

The provider redacts token and secret-like values before including API error messages in Terraform diagnostics.

## Request Example

```json
{
  "name": "ledger-studio",
  "product_code": "bankport",
  "redirect_uris": ["https://ledger.example.com/oauth/callback"],
  "scopes": ["accounts:read", "payments:write"],
  "status": "active"
}
```

## Response Example

```json
{
  "id": "app_0001",
  "name": "ledger-studio",
  "product_code": "bankport",
  "redirect_uris": ["https://ledger.example.com/oauth/callback"],
  "scopes": ["accounts:read", "payments:write"],
  "status": "active",
  "client_id": "client_app_0001",
  "client_secret": "client_secret_app_0001"
}
```

`client_secret` is returned by the fake API for test determinism and is marked sensitive in Terraform schema.
