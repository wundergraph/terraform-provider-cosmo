variable "monograph_name" {
  default = "tf-guide-cosmo-monograph-contract"
}

variable "monograph_namespace" {
  default = "tf-guide-cosmo-monograph-contract"
}

variable "monograph_graph_url" {
  default = "http://example.com/graphql"
}

variable "monograph_routing_url" {
  default = "http://example.com/routing"
}

variable "contract_name" {
  type    = string
  default = "test"
}

variable "contract_namespace" {
  type    = string
  default = "default"
}

variable "contract_routing_url" {
  type    = string
  default = "http://example.com/routing"
}

variable "contract_exclude_tags" {
  type    = list(string)
  default = []
}

variable "api_url" {
  type    = string
  default = "http://example.com/graphql"
}

variable "api_key" {
  type = string
}