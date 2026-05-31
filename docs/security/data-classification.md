# Data Classification

| Data | Classification | Storage | Handling |
| --- | --- | --- | --- |
| Provider endpoint | Internal | HCL/state | Not sensitive. |
| Provider token | Secret | Terraform configuration/environment | Marked sensitive in schema; prefer env or secret store. |
| Client ID | Internal identifier | Terraform state | Not secret. |
| Client secret | Secret | Terraform state | Sensitive schema; encrypted backend required. |
| Webhook signing secret | Secret | Terraform state | Sensitive schema; rotate on compromise. |
| Sandbox API key token | Secret | Terraform state | Sensitive schema; restrict outputs. |
| Rate-limit values | Internal | Terraform state | Review for abuse impact. |
