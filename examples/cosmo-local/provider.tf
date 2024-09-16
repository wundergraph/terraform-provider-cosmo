terraform {
  required_providers {
    minikube = {
      source  = "scott-the-programmer/minikube"
      version = "0.4.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "2.32.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "2.15.0"
    }
    cosmo = {
      source  = "terraform.local/wundergraph/cosmo"
      version = "0.0.1"
    }
    time = {
      source  = "hashicorp/time"
      version = "0.12.1"
    }
  }
}

provider "kubernetes" {
  host                   = module.minikube.host
  client_certificate     = module.minikube.client_certificate
  client_key             = module.minikube.client_key
  cluster_ca_certificate = module.minikube.cluster_ca_certificate
}

provider "helm" {
  kubernetes {
    host                   = module.minikube.host
    client_certificate     = module.minikube.client_certificate
    client_key             = module.minikube.client_key
    cluster_ca_certificate = module.minikube.cluster_ca_certificate
  }
}

provider "cosmo" {
  api_url = var.api_url
  api_key = var.api_key
}

