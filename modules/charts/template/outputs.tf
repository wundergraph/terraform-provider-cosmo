output "manifests" {
  value = data.helm_template.this.manifests
}

output "manifest" {
  value = data.helm_template.this.manifest
}

output "notes" {
  value = data.helm_template.this.notes
}