resource "cosmo_subgraph" "test" {
  name        = var.name
  namespace   = var.namespace
  routing_url = var.routing_url
  labels      = var.labels
  schema      = var.schema
}