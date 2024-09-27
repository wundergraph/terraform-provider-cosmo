module "cosmo_contract" {
  # only for test purposes give each federated graph a contract
  for_each = var.federated_graphs
  source   = "../../resources/cosmo_contract"

  name              = "${each.value.name}-contract"
  namespace         = module.cosmo_federated_graph[each.key].namespace_name
  source_graph_name = module.cosmo_federated_graph[each.key].federated_graph_name
  routing_url       = "http://localhost:3000"
  exclude_tags      = ["backend"]
}

