# Glossary

| Term | Meaning |
| --- | --- |
| API Product | BankPort, PixGuard, or SettleFlow capability exposed through a platform API. |
| Partner App | OAuth-style application representing an external partner integration. |
| Webhook Endpoint | HTTPS destination that receives platform events for a partner app. |
| Rate-Limit Policy | Remote control-plane rule that limits API calls by subject. |
| Sandbox Environment | Test environment provisioned for a partner and product set. |
| Terraform State | Cached desired/observed provider state, including sensitive values. |
| Drift | Difference between Terraform state/configuration and the remote API. |
| Import | Bringing an existing remote resource ID under Terraform management. |
| Generated Secret | Secret returned by the API during create or rotation. |
