terraform {
  required_providers {
    cosmo = {
      source  = "wundergraph/cosmo"
      version = "0.0.1"
    }
  }
}

data "cosmo_subgraph" "test" {
  name      = var.name
  namespace = var.namespace
}