resource "cosmo_namespace" "namespace" {
  name = var.namespace
}

resource "cosmo_federated_graph" "federated_graph" {
  name           = var.federated_graph.name
  routing_url    = var.federated_graph.routing_url
  namespace      = cosmo_namespace.namespace.name
  label_matchers = var.federated_graph.label_matchers
  depends_on     = [cosmo_subgraph.subgraph]
}

resource "cosmo_subgraph" "subgraph" {
  for_each = var.subgraphs

  name        = each.value.name
  namespace   = cosmo_namespace.namespace.name
  routing_url = each.value.routing_url
  schema      = each.value.schema
  labels      = each.value.labels
  readme      = each.value.readme
}

resource "cosmo_router_token" "router_token" {
  name       = var.router_token_name
  namespace  = cosmo_namespace.namespace.name
  graph_name = cosmo_federated_graph.federated_graph.name
}