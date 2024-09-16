resource "helm_release" "this" {
  name       = var.release_name
  namespace  = var.chart.namespace
  repository = var.chart.repository
  version    = var.chart.version
  chart      = var.chart.name
  values     = concat([file(var.chart.init_values)], var.chart.values)

  wait          = true
  wait_for_jobs = true

  cleanup_on_fail = false
  atomic          = false
  verify          = false

  dynamic "set" {
    for_each = var.chart.set
    content {
      name  = set.key
      value = set.value
    }
  }
}