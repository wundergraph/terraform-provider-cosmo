terraform {
  required_providers {
    carvel = {
      source  = "carvel-dev/carvel"
      version = "0.11.2"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "2.15.0"
    }
  }
}