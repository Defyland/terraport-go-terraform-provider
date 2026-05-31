# Authorization Matrix

| Operation | Required Platform Capability | Provider Behavior | Test Evidence |
| --- | --- | --- | --- |
| Read API product | `products:read` | `GET /v1/products/{code}` | `TestAccAPIProductDataSource` |
| Create partner app | `partner_apps:write` | `POST /v1/partner-apps` | `TestAccPartnerAppLifecycleImportDrift` |
| Rotate app secret | `partner_apps:rotate_secret` | `POST /rotate-secret` when version increases | `TestAccPartnerAppLifecycleImportDrift` |
| Create webhook | `webhooks:write` | `POST /v1/webhook-endpoints` | `TestAccWebhookEndpointLifecycleImport` |
| Create rate policy | `rate_limits:write` | `POST /v1/rate-limit-policies` | `TestAccRateLimitPolicyLifecycleImport` |
| Create sandbox | `sandboxes:write` | `POST /v1/sandbox-environments` | `TestAccSandboxEnvironmentLifecycleImport` |

The fake API only validates bearer token presence/value. Real tenant authorization is a remote API responsibility and is documented as out of scope for this repository.
