locals {
  cosmo_router = {
    release_name = var.cosmo_router.release_name
    chart = {
      name        = var.cosmo_router.chart.name
      version     = var.cosmo_router.chart.version
      namespace   = var.cosmo_router.chart.namespace
      repository  = var.cosmo_router.chart.repository
      values      = concat(var.cosmo_router.chart.values, [])
      init_values = var.cosmo_router.chart.init_values
      set = merge({
        "image.registry"                           = "ghcr.io"
        "image.repository"                         = "wundergraph/cosmo/router"
        "configuration.graphApiToken"              = module.cosmo_federated_graph.router_token
        "configuration.controlplaneUrl"            = "http://cosmo-controlplane.cosmo.svc.cluster.local:3001"
        "configuration.cdnUrl"                     = "http://cosmo-cdn.cosmo.svc.cluster.local:8787"
        "configuration.otelCollectorUrl"           = "http://cosmo-otelcollector.cosmo.svc.cluster.local:4318"
        "configuration.graphqlMetricsCollectorUrl" = "http://cosmo-graphqlmetrics.cosmo.svc.cluster.local:4005"
      }, var.cosmo_router.chart.set)
    }
  }
}

module "minikube" {
  source             = "../minikube"
  minikube_clusters  = var.minikube
  kubernetes_version = var.kubernetes_version
}

resource "time_sleep" "wait_for_minikube" {
  create_duration = "10s"

  depends_on = [module.minikube]
}

module "cosmo_charts" {
  source = "../charts/template"
  chart  = var.cosmo.chart

  depends_on = [time_sleep.wait_for_minikube]
}

resource "kubernetes_namespace" "cosmo_namespace" {
  metadata {
    name = var.cosmo.chart.namespace
  }

  depends_on = [time_sleep.wait_for_minikube]
}

module "cosmo_release" {
  source = "../charts/release"
  chart  = var.cosmo.chart

  release_name = var.cosmo.release_name

  depends_on = [time_sleep.wait_for_minikube]
}

module "cosmo_router_release" {
  source = "../charts/release"
  chart  = local.cosmo_router.chart

  release_name = local.cosmo_router.release_name

  depends_on = [
    module.cosmo_release,
    module.cosmo_federated_graph
  ]
}