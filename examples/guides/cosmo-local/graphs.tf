// 1. Create a random prefix for the module
// this prefixes the test resources to avoid conflicts
// and to find them in the cosmo ui easier
resource "random_string" "module_prefix" {
  length  = 6
  special = false
}

// 2. Lowercase the prefix and store it in locals
// namespace have to be lowercase therefore we lowercase the prefix
locals {
  prefix = lower(random_string.module_prefix.result)
  schema = <<EOF
type Wunder {
  title: String
}

type Graph {
  name: String
  wunders: [Wunder]
}
EOF
}

// 3. Create the federated graph
// this module wraps generating a federated graph and related subgraphs
// the resources are deployed within the given namespace
module "cosmo_federated_graph" {
  for_each          = var.federated_graphs
  source            = "../../../modules/cosmo-federated-graph"
  namespace         = "${each.key}-${var.stage}-${local.prefix}"
  router_token_name = "${each.key}-${var.stage}-${local.prefix}-router-token"

  // 3.1. The federated graph configuration
  // this represents the federated graph in cosmo
  // see docs/resources/federated_graph.md for more information
  federated_graph = {
    name        = "${each.key}-${var.stage}-${local.prefix}-federated-graph"
    readme      = "The federated graph for the local setup"
    routing_url = "http://localhost:3000"
    label_matchers = [
      "team=backend",
      "stage=${var.stage}",
      "graph=${each.key}"
    ]
  }

  // 3.2. The subgraphs to be added to the federated graph
  // this represents the subgraphs for the federated graph
  // see docs/resources/subgraph.md for more information
  subgraphs = {
    "spacex" = {
      name        = "${each.key}-${var.stage}-${local.prefix}-spacex"
      routing_url = "https://spacex-production.up.railway.app/graphql"
      schema      = var.switch_schema ? local.schema : file("${path.module}/schema/spacex-api.graphql")
      readme      = <<EOF
# Overview 

SpaceX is a company that builds spacecraft and rockets.
      EOF
      labels = {
        "team"  = "backend"
        "stage" = var.stage
        "graph" = each.key
      }
    }
  }

  depends_on = [
    module.minikube,
    module.cosmo_release
  ]
}