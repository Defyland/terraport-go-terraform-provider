variable "terraport_endpoint" {
  type        = string
  description = "BankPort-compatible API endpoint."
}

variable "terraport_token" {
  type        = string
  sensitive   = true
  description = "Bearer token for the BankPort-compatible API."
}

provider "terraport" {
  endpoint             = var.terraport_endpoint
  token                = var.terraport_token
  timeout_ms           = 10000
  retry_max_attempts   = 4
  retry_min_delay_ms   = 100
}

resource "terraport_bankport_partner_app" "ledger_studio" {
  name                  = "ledger-studio"
  product_code          = "bankport"
  redirect_uris         = ["https://ledger.example.com/oauth/callback"]
  scopes                = ["accounts:read", "payments:write", "webhooks:write"]
  client_secret_version = 1
}

output "client_id" {
  value = terraport_bankport_partner_app.ledger_studio.client_id
}
