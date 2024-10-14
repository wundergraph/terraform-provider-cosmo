locals {

}

// 1. Install minikube on which cosmo will be deployed
module "minikube" {
  source             = "../../../modules/minikube"
  minikube_clusters  = var.minikube
  kubernetes_version = var.kubernetes_version
}

// 2. Wait for minikube to be ready to avoid race conditions with helm
resource "time_sleep" "wait_for_minikube" {
  create_duration = "30s"

  depends_on = [module.minikube]
}

// 3. Render the cosmo charts (used for local development and debugging)
module "cosmo_charts" {
  source = "../../../modules/charts/template"
  chart  = var.cosmo.chart

  depends_on = [time_sleep.wait_for_minikube]
}

// 4. Install the cosmo helm release 
// see var.cosmo.release_name and var.cosmo.chart for more details
module "cosmo_release" {
  source = "../../../modules/charts/release"
  chart  = var.cosmo.chart

  release_name = var.cosmo.release_name

  depends_on = [time_sleep.wait_for_minikube]
}

// 5. Install the cosmo router helm release
// see local.cosmo_router.release_name and local.cosmo_router.chart for more details
// this happens after graphs.tf was applied after the router token was created
module "cosmo_router_release" {
  for_each = var.federated_graphs
  source   = "../../../modules/charts/release"

  chart = {
    name             = var.cosmo_router.chart.name
    version          = var.cosmo_router.chart.version
    namespace        = var.cosmo_router.chart.namespace
    repository       = var.cosmo_router.chart.repository
    values           = concat(var.cosmo_router.chart.values, [])
    init_values      = var.cosmo_router.chart.init_values
    create_namespace = false
    set = merge({
      "configuration.graphApiToken"              = module.cosmo_federated_graph[each.key].router_token
      "configuration.controlplaneUrl"            = "http://cosmo-controlplane.${var.cosmo.chart.namespace}.svc.cluster.local:3001"
      "configuration.cdnUrl"                     = "http://cosmo-cdn.${var.cosmo.chart.namespace}.svc.cluster.local:8787"
      "configuration.otelCollectorUrl"           = "http://cosmo-otelcollector.${var.cosmo.chart.namespace}.svc.cluster.local:4318"
      "configuration.graphqlMetricsCollectorUrl" = "http://cosmo-graphqlmetrics.${var.cosmo.chart.namespace}.svc.cluster.local:4005"
    }, var.cosmo_router.chart.set)
  }

  release_name = "${each.key}-${var.stage}-${local.prefix}-router"

  depends_on = [
    module.cosmo_release,
    module.cosmo_federated_graph
  ]
}