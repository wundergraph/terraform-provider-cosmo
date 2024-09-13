terraform {
  required_providers {
    cosmo = {
      source  = "terraform.local/wundergraph/cosmo"
      version = "0.0.1"
    }
  }
}

provider "cosmo" {
  api_url = "cosmo_669b576aaadc10ee1ae81d9193425705"
  api_key = "http://localhost:3001"
}
