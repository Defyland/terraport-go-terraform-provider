# Observability

Terraport is not a long-running HTTP service, so it does not expose `/metrics`, health checks, or readiness endpoints directly. The relevant observability surface is Terraform diagnostics plus remote BankPort API telemetry.

## Implemented Signals

- Client request counter.
- Retry counter.
- Rate-limit response counter.
- Server-error response counter.
- Fake API method/path counters used by tests.
- Sanitized Terraform diagnostics for API errors.

## Expected Platform Signals

A real BankPort API should expose:

- `bankport_control_plane_requests_total`
- `bankport_control_plane_rate_limited_total`
- `bankport_partner_app_mutations_total`
- `bankport_secret_rotations_total`
- request and correlation IDs for Terraform apply windows

## Alerts

- Sustained `429` responses during Terraform apply windows.
- Increased `401` responses from CI identities.
- Secret rotation failures.
- Drift-related manual changes in platform audit logs.

## Dashboard

The provider does not run a Grafana dashboard itself. A platform dashboard should combine BankPort API request metrics with Terraform Cloud or CI apply outcomes.
