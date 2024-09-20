module "cosmo_monograph" {
  source = "../../resources/cosmo_monograph"

  monograph_name        = var.monograph_name
  monograph_namespace   = var.monograph_namespace
  monograph_graph_url   = var.monograph_graph_url
  monograph_routing_url = var.monograph_routing_url
}