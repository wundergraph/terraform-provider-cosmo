resource "cosmo_monograph" "example" {
  name        = var.monograph_name
  namespace   = var.monograph_namespace
  graph_url   = var.monograph_graph_url
  routing_url = var.monograph_routing_url
  schema      = var.monograph_schema
}