terraform {
  required_providers {
    cosmo = {
      source  = "wundergraph/cosmo"
      version = "0.0.1"
    }
  }
}

data "cosmo_monograph" "example" {
  name = var.name
}