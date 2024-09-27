// your stages represented by different namespaces
locals {
  stages = {
    dev  = {},
    stg  = {},
    prod = {}
  }
}

// your subgraphs, which are deployed to each stage
locals {
  subgraphs = {
    "product-api" = {
      routing_url = "http://product-api:3000/graphql"
      labels = {
        "team" = "backend"
      }
    },
    "employees-api" = {
      routing_url = "http://employees-api:3000/graphql"
      labels = {
        "team" = "backend"
      }
    },
    "family-api" = {
      routing_url = "http://family-api:3000/graphql"
      labels = {
        "team" = "backend"
      }
    },
    "hobbies-api" = {
      routing_url = "http://hobbies-api:3000/graphql"
      labels = {
        "team" = "backend"
      }
    },
    "availability-api" = {
      routing_url = "http://availability-api:3000/graphql"
      labels = {
        "team" = "backend"
      }
    }
  }

  // Helper, used to make the subgraphs above staged
  // {
  //   "dev-product-api" = {
  //     "stage" = "dev"
  //     "subgraph" = "product-api"
  //     "routing_url" = "http://product-api:3000/graphql"
  //     "labels" = {
  //       "team" = "backend"
  //     }
  //   }
  // }
  stage_subgrahs = merge(flatten([
    for key, value in local.stages : {
      for subgraph, subgraph_value in local.subgraphs :
      "${key}-${subgraph}" => {
        "stage"       = key
        "subgraph"    = subgraph
        "routing_url" = subgraph_value.routing_url
        "labels"      = subgraph_value.labels
      }
  }])...)
}

// create a namespace for each stage
// dev-namespace, stg-namespace, prod-namespace
resource "cosmo_namespace" "namespace" {
  for_each = local.stages

  name = "${each.key}-namespace"
}

// create a federated graph for each stage
// dev-federated-graph, stg-federated-graph, prod-federated-graph
resource "cosmo_federated_graph" "federated_graph" {
  for_each = local.stages

  name        = "${each.key}-federated-graph"
  routing_url = "http://${each.key}.localhost:3000"
  namespace   = cosmo_namespace.namespace[each.key].name

  label_matchers = ["team=backend", "stage=${each.key}"]

  depends_on = [cosmo_subgraph.subgraph]
}

// create each stages subgraph
resource "cosmo_subgraph" "subgraph" {
  for_each = local.stage_subgrahs

  name      = "${each.value.subgraph}-${each.value.stage}-subgraph"
  namespace = cosmo_namespace.namespace[each.value.stage].name

  routing_url = each.value.routing_url

  // merge graph labels with the stage label
  labels = merge(each.value.labels, {
    "stage" = each.value.stage
  })
}

// create a router token for each stage
resource "cosmo_router_token" "router_token" {
  for_each = local.stages

  name       = "${each.key}-router-token"
  namespace  = cosmo_namespace.namespace[each.key].name
  graph_name = cosmo_federated_graph.federated_graph[each.key].name
}

output "dev_router_token" {
  value     = cosmo_router_token.router_token["dev"].token
  sensitive = true
}

output "stg_router_token" {
  value     = cosmo_router_token.router_token["stg"].token
  sensitive = true
}

output "prod_router_token" {
  value     = cosmo_router_token.router_token["prod"].token
  sensitive = true
}

// used to debug sensitive values
// resource "local_file" "router_tokens" {
//   content = jsonencode(cosmo_router_token.router_token)
//   filename = "router_tokens.json"
// }