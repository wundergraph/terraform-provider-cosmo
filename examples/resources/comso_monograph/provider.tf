terraform {
  required_providers {
    cosmo = {
      source  = "wundergraph/cosmo"
      version = "0.0.1"
    }
  }
}

provider "cosmo" {
  api_url = var.api_url
  api_key = var.api_key
}