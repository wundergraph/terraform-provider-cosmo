terraform {
  required_providers {
    cosmo = {
      source  = "terraform.local/wundergraph/cosmo"
      version = "0.0.1"
    }
  }
}

provider "cosmo" {
  api_url = var.api_url
  api_key = var.api_key
}
