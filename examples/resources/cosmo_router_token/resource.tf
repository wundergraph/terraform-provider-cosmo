terraform {
  required_providers {
    cosmo = {
      source  = "wundergraph/cosmo"
      version = "0.0.1"
    }
  }
}

resource "cosmo_router_token" "test" {
  name       = var.name
  graph_name = var.graph_name
  namespace  = var.namespace
}