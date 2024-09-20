resource "cosmo_contract" "test" {
  name         = var.name
  namespace    = var.namespace
  source       = var.source_graph_name
  routing_url  = var.routing_url
  exclude_tags = var.exclude_tags
}
