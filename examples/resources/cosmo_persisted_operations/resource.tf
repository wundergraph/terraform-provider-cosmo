resource "cosmo_persisted_operations" "test" {
  federated_graph_name = var.federated_graph_name
  namespace            = var.namespace
  client_name          = var.client_name

  operations = {
    for file in fileset(path.module, "operations/*.graphql") :
    trimsuffix(basename(file), ".graphql") => file("${path.module}/${file}")
  }
}
