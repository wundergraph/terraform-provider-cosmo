output "host" {
  value = minikube_cluster.this.host
}

output "client_certificate" {
  value = minikube_cluster.this.client_certificate
}

output "client_key" {
  value = minikube_cluster.this.client_key
}

output "cluster_ca_certificate" {
  value = minikube_cluster.this.cluster_ca_certificate
}
