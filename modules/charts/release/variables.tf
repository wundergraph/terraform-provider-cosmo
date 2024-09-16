variable "chart" {
  type = object({
    name        = string
    version     = string
    namespace   = string
    repository  = string
    values      = list(string)
    init_values = string
    set         = map(string)
  })
  default = {
    name        = "cosmo"
    version     = "0.11.1"
    namespace   = "cosmo"
    repository  = "oci://ghcr.io/wundergraph/cosmo/helm-charts"
    values      = []
    init_values = ""
    set         = {}
  }
}

variable "release_name" {
  type    = string
  default = "cosmo"
}