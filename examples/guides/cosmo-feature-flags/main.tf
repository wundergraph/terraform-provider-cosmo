module "cosmo_namespace" {
  source = "../../resources/cosmo_namespace"

  name = var.namespace
}

module "cosmo_federated_graph" {
  source = "../../resources/cosmo_federated_graph"

  name           = var.federated_graph_name
  namespace      = module.cosmo_namespace.name
  routing_url    = var.federated_graph_routing_url
  label_matchers = var.federated_graph_label_matchers
}

module "cosmo_subgraph" {
  source = "../../resources/cosmo_subgraph"

  name        = var.subgraph_name
  namespace   = module.cosmo_namespace.name
  routing_url = var.subgraph_routing_url
  schema      = var.subgraph_schema
  labels      = var.subgraph_labels
}

module "cosmo_feature_subgraph" {
  source = "../../resources/cosmo_feature_subgraph"

  name                  = var.feature_subgraph_name
  namespace             = module.cosmo_namespace.name
  base_subgraph_name    = module.cosmo_subgraph.name
  readme                = var.feature_subgraph_readme
  routing_url           = var.feature_subgraph_routing_url
  schema                = var.feature_subgraph_schema
  subscription_protocol = var.subscription_protocol
  subscription_url      = var.subscription_url
  websocket_subprotocol = var.websocket_subprotocol
}

module "cosmo_feature_flag" {
  source = "../../resources/cosmo_feature_flag"

  name              = var.feature_flag_name
  namespace         = module.cosmo_namespace.name
  feature_subgraphs = [module.cosmo_feature_subgraph.name]
  is_enabled        = var.feature_flag_is_enabled
  labels            = var.feature_flag_labels
}