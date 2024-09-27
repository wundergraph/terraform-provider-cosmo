locals {
  minikube_ip = replace(replace(module.minikube.host, "https://", ""), ":8443", "")
}

output "hosts" {
  value = <<EOF
    # WunderGraph
    ${local.minikube_ip} studio.wundergraph.local
    ${local.minikube_ip} controlplane.wundergraph.local
    ${local.minikube_ip} router.wundergraph.local
    ${local.minikube_ip} keycloak.wundergraph.local
    ${local.minikube_ip} otelcollector.wundergraph.local
    ${local.minikube_ip} graphqlmetrics.wundergraph.local
    ${local.minikube_ip} cdn.wundergraph.local
    EOF
}