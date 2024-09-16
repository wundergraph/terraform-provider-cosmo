resource "local_file" "deployment" {
  for_each = var.write_deployment_yaml ? {
    "deployment" = module.cosmo_charts.manifest
  } : {}

  content  = module.cosmo_charts.manifest
  filename = "${path.module}/${var.stage}-deployment.yaml"
}