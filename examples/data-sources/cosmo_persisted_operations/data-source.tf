data "cosmo_persisted_operations" "test" {
  federated_graph_name = var.federated_graph_name
  namespace            = var.namespace
  client_name          = var.client_name
}
