variable "stage" {
  type        = string
  description = "The stage of the deployment"
  default     = "local"
}

variable "write_deployment_yaml" {
  type        = bool
  description = "Whether to write the deployment yaml to a file"
  default     = false
}

variable "kubernetes_version" {
  type    = string
  default = "v1.30.0"
}

variable "minikube" {
  type = object({
    cluster_name = string
    driver       = string
    addons       = list(string)
  })
  default = {
    cluster_name = "cosmo-local-docker"
    driver       = "docker"
    addons = [
      "dashboard",
      "default-storageclass",
      "ingress",
      "storage-provisioner"
    ]
  }
}

variable "cosmo" {
  type = object({
    release_name = string
    chart = object({
      name             = string
      version          = string
      namespace        = string
      repository       = string
      values           = list(string)
      init_values      = string
      set              = map(string)
      create_namespace = bool
    })
  })
  description = "The cosmo chart to deploy"
  default = {
    release_name = "cosmo"
    chart = {
      name             = "cosmo"
      version          = "0.12.2"
      namespace        = "cosmo"
      repository       = "oci://ghcr.io/wundergraph/cosmo/helm-charts"
      values           = []
      init_values      = "./values/cosmo-values.yaml"
      set              = {}
      create_namespace = true
    }
  }
}

variable "cosmo_router" {
  type = object({
    release_name = string
    chart = object({
      name        = string
      version     = string
      namespace   = string
      repository  = string
      values      = list(string)
      init_values = string
      set         = map(string)
    })
  })
  description = "The cosmo router chart to deploy"
  default = {
    release_name = "cosmo-router"
    chart = {
      name        = "router"
      version     = "0.8.0"
      namespace   = "cosmo"
      repository  = "oci://ghcr.io/wundergraph/cosmo/helm-charts"
      values      = []
      init_values = "./values/router-values.yaml"
      set = {
        "image.registry"   = "ghcr.io"
        "image.repository" = "wundergraph/cosmo/router"
        "image.version"    = "0.110.1"
      }
    }
  }
}

variable "api_url" {
  type        = string
  description = "The helm Charts control plane url"
}

variable "api_key" {
  type        = string
  description = "The helm Charts control plane api key"
}

variable "switch_schema" {
  type        = bool
  description = "Switch the schema to use for the federated graph"
  default     = false
}

variable "federated_graphs" {
  type = map(object({
    name = string
  }))
  description = "Switch the schema to use for the federated graph"
  default = {
    "spacex" = {
      name = "spacex"
    }
  }
}
