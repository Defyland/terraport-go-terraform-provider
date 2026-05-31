variable "terraport_endpoint" {
  type = string
}

variable "terraport_token" {
  type      = string
  sensitive = true
}

provider "terraport" {
  endpoint   = var.terraport_endpoint
  token      = var.terraport_token
  timeout_ms = 15000
}

data "terraport_bankport_api_product" "bankport" {
  product_code = "bankport"
}

resource "terraport_bankport_sandbox_environment" "partner_lab" {
  name     = "partner-lab"
  products = ["bankport", "pixguard", "settleflow"]
  region   = "sa-east-1"
}

output "bankport_product_name" {
  value = data.terraport_bankport_api_product.bankport.name
}
