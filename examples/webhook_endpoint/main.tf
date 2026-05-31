variable "terraport_endpoint" {
  type = string
}

variable "terraport_token" {
  type      = string
  sensitive = true
}

provider "terraport" {
  endpoint = var.terraport_endpoint
  token    = var.terraport_token
}

resource "terraport_bankport_partner_app" "platform_ops" {
  name                  = "platform-ops"
  product_code          = "bankport"
  redirect_uris         = ["https://ops.example.com/oauth/callback"]
  scopes                = ["webhooks:write"]
  client_secret_version = 1
}

resource "terraport_bankport_webhook_endpoint" "settlements" {
  partner_app_id         = terraport_bankport_partner_app.platform_ops.id
  url                    = "https://ops.example.com/webhooks/bankport"
  event_types            = ["payment.settled", "partner_app.secret_rotated"]
  enabled                = true
  signing_secret_version = 1
}
