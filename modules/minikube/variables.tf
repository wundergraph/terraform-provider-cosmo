variable "minikube_clusters" {
  type = object({
    driver       = string
    cluster_name = string
    nodes        = optional(number)
    addons       = list(string)
  })
  default = {
    driver       = "docker"
    cluster_name = "cosmo-local-docker"
    nodes        = 1
    addons       = ["dashboard", "default-storageclass", "ingress", "storage-provisioner"]
  }
}

variable "kubernetes_version" {
  type    = string
  default = "v1.30.0"
}