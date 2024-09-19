resource "cosmo_federated_graph" "test" {
  name           = var.name
  routing_url    = var.routing_url
  namespace      = var.namespace
  label_matchers = var.label_matchers
}
