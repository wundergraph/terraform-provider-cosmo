resource "random_string" "module_prefix" {
  length  = 6
  special = false
}

locals {
  prefix = lower(random_string.module_prefix.result)
}

module "cosmo_federated_graph" {
  source            = "../cosmo-federated-graph"
  namespace         = "${var.stage}-${local.prefix}"
  router_token_name = "${var.stage}-${local.prefix}-router-token"

  federated_graph = {
    name        = "${local.prefix}-federated-graph"
    readme      = "The federated graph for the local setup"
    routing_url = "http://localhost:3000"
    label_matchers = [
      "team=backend",
    ]
  }
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