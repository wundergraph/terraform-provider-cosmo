resource "minikube_cluster" "this" {
  driver       = var.minikube_clusters.driver
  cluster_name = var.minikube_clusters.cluster_name
  addons       = var.minikube_clusters.addons
  wait         = ["all"]
}
