terraform {
  required_providers {
    cosmo = {
      source  = "terraform.local/wundergraph/cosmo"
      version = "0.0.1"
    }
  }
}

provider "cosmo" {
  api_url = "http://localhost:3001"
  api_key = "cosmo_669b576aaadc10ee1ae81d9193425705"
}
