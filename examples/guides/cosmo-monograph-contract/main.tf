module "cosmo_namespace" {
  source = "../../resources/cosmo_namespace"

  name = var.monograph_namespace
}

module "cosmo_monograph" {
  source = "../../resources/cosmo_monograph"

  monograph_name        = var.monograph_name
  monograph_namespace   = module.cosmo_namespace.name
  monograph_graph_url   = var.monograph_graph_url
  monograph_routing_url = var.monograph_routing_url
  monograph_schema = var.monograph_schema
}

module "cosmo_contract" {
  source = "../../resources/cosmo_contract"

  name              = var.contract_name
  namespace         = module.cosmo_namespace.name
  source_graph_name = module.cosmo_monograph.name
  routing_url       = var.contract_routing_url
  exclude_tags      = var.contract_exclude_tags
}
