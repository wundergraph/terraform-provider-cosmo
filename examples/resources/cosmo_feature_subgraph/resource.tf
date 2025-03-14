resource "cosmo_feature_subgraph" "example" {
  name                  = var.name
  namespace             = var.namespace
  routing_url           = var.routing_url
  subscription_url      = var.subscription_url
  subscription_protocol = var.subscription_protocol
  websocket_subprotocol = var.websocket_subprotocol
  base_subgraph_name    = var.base_subgraph_name
  schema                = var.schema
  readme                = var.readme
}
