module "resource_cosmo_namespace" {
  source = "../resources/cosmo_namespace"

  name = "terraform-namespace-demo"
}

module "resource_cosmo_subgraph" {
  source = "../resources/cosmo_subgraph"

  name        = "subgraph-1"
  namespace   = module.resource_cosmo_namespace.name
  routing_url = "http://example.com/routing"
  schema      = file("schema.graphql")
  labels = {
    "team"  = "backend"
    "stage" = "dev"
  }
}

module "resource_cosmo_federated_graph" {
  source = "../resources/cosmo_federated_graph"

  name           = "terraform-federated-graph-demo"
  routing_url    = "http://localhost:3000"
  namespace      = module.resource_cosmo_namespace.name
  label_matchers = ["team=backend"]

  depends_on = [module.resource_cosmo_subgraph]
}

module "resource_cosmo_contract" {
  source = "../resources/cosmo_contract"

  name              = "terraform-contract-demo"
  namespace         = module.resource_cosmo_namespace.name
  routing_url       = module.resource_cosmo_federated_graph.routing_url
  source_graph_name = module.resource_cosmo_federated_graph.name
}

module "data_cosmo_federated_graph" {
  source = "../data-sources/cosmo_federated_graph"

  name      = module.resource_cosmo_federated_graph.name
  namespace = module.resource_cosmo_namespace.name

  // This is necessary, as ID is computed, but the datasource depends on the not computed name. 
  // Only needed when creation and reading happen in the same apply.
  depends_on = [module.resource_cosmo_federated_graph]
}