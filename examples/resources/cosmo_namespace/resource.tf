terraform {
  required_providers {
    cosmo = {
      source  = "wundergraph/cosmo"
      version = "0.0.1"
    }
  }
}

resource "cosmo_namespace" "test" {
  name = var.name
}