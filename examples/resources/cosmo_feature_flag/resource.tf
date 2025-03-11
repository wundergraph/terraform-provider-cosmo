resource "cosmo_feature_flag" "example" {
  name              = var.name
  namespace         = var.namespace
  feature_subgraphs = var.feature_subgraphs
  labels            = var.labels
  is_enabled        = var.is_enabled
}