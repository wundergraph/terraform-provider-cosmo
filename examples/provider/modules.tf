resource "random_string" "module_prefix" {
  length  = 6
  special = false
}

locals {
  prefix = lower(random_string.module_prefix.result)
}

module "cosmo-federated-graph" {
  source = "../../modules/cosmo-federated-graph"
  namespace = "${local.prefix}-cosmo-module"
  router_token_name = "${local.prefix}-router-token"

  federated_graph = {
    name        = "${local.prefix}-federated-graph"
    routing_url = "http://localhost:3000"
    label_matchers = [
      "team=backend",
      "stage=dev"
    ]
  }
  subgraphs = {
    "subgraph-1" = {
      name        = "${local.prefix}-subgraph-1"
      routing_url = "http://example.com/routing"
      labels = {
        "team"  = "backend"
        "stage" = "dev"
      }
    }
  }
}

