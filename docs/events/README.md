# Events

Terraport does not publish messages or run a broker. In a real BankPort control plane, provider operations would create audit events. The fake API does not persist an outbox, but the event contract below documents what downstream audit consumers would expect.

## Envelope

```json
{
  "event_id": "evt_01J...",
  "event_type": "bankport.partner_app.created",
  "schema_version": "1",
  "occurred_at": "2026-05-31T12:00:00Z",
  "producer": "bankport-platform-api",
  "correlation_id": "tf-apply-123",
  "resource_id": "app_0001"
}
```

## Compatibility Rules

- Do not remove required fields without a new schema version.
- Additive fields must be optional.
- Consumers must handle duplicate delivery.
- Consumers must use `event_id` for idempotency and `correlation_id` for trace stitching.

## Consumer Expectations

Audit consumers should never receive raw `client_secret`, `signing_secret`, or `api_key_token` values. Secret rotation events may include a secret version but not the secret.
