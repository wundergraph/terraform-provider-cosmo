terraform {
  required_providers {
    cosmo = {
      source  = "wundergraph/cosmo"
      version = "0.0.1"
    }
  }
}

resource "cosmo_subgraph" "test" {
  name        = var.name
  namespace   = var.namespace
  routing_url = var.routing_url
  labels      = var.labels
}