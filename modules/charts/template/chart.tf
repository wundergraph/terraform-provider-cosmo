data "helm_template" "this" {
  name       = var.chart.name
  namespace  = var.chart.namespace
  repository = var.chart.repository
  chart      = var.chart.name
  version    = var.chart.version
  values     = concat([file(var.chart.init_values)], var.chart.values)

  dynamic "set" {
    for_each = var.chart.set
    content {
      name  = set.key
      value = set.value
    }
  }
}
