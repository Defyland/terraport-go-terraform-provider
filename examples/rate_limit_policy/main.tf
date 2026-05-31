variable "terraport_endpoint" {
  type = string
}

variable "terraport_token" {
  type      = string
  sensitive = true
}

provider "terraport" {
  endpoint           = var.terraport_endpoint
  token              = var.terraport_token
  retry_max_attempts = 4
}

resource "terraport_bankport_rate_limit_policy" "partner_burst" {
  product_code        = "bankport"
  subject_type        = "partner_app"
  subject_id          = "app_existing_123"
  requests_per_minute = 1200
  burst_limit         = 120
  mode                = "enforce"
}
