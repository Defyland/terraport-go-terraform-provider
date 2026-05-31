# Roadmap

## Now

- Provider configuration through block and environment variables.
- CRUD, import, sensitive fields, drift tests, and fake API acceptance tests.
- Examples for partner apps, webhooks, rate limits, and sandbox environments.

## Next

- Add generated provider documentation from schemas.
- Add real pagination and correlation ID propagation once the platform API shape is stable.
- Add PixGuard and SettleFlow resource families when product control-plane APIs are defined.

## Later

- Publish to a private Terraform Registry.
- Add provider upgrade tests for schema migrations.
- Add batched refresh endpoints if large Terraform workspaces start hitting remote rate limits.
