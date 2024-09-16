// 1. Create a random prefix for the module
resource "random_string" "module_prefix" {
  length  = 6
  special = false
}

// 2. Lowercase the prefix and store it in locals
locals {
  prefix = lower(random_string.module_prefix.result)
}

// 3. Create the federated graph
// this creates a federated graph in cosmo given a set of subgraphs
// and the namespace to deploy the graph to
module "cosmo_federated_graph" {
  source            = "../../modules/cosmo-federated-graph"
  namespace         = "${var.stage}-${local.prefix}"
  router_token_name = "${var.stage}-${local.prefix}-router-token"

  // 3.1. The federated graph configuration
  federated_graph = {
    name        = "${var.stage}-${local.prefix}-federated-graph"
    readme      = "The federated graph for the local setup"
    routing_url = "http://localhost:3000"
    label_matchers = [
      "team=backend",
      "stage=${var.stage}"
    ]
  }

  // 3.2. The subgraphs to be added to the federated graph
  subgraphs = {
    "spacex" = {
      name        = "${var.stage}-${local.prefix}-spacex"
      routing_url = "https://spacex-production.up.railway.app/graphql"
      schema      = file("${path.module}/schema/spacex-api.graphql")
      readme      = <<EOF
# Overview 

SpaceX is a company that builds spacecraft and rockets.
      EOF
      labels = {
        "team"  = "backend"
        "stage" = "${var.stage}"
      }
    }
  }

  depends_on = [
    module.minikube,
    module.cosmo_release
  ]
}