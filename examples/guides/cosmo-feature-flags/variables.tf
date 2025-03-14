variable "namespace" {
  default = "tf-guide-cosmo-feature-flags"
}

variable "federated_graph_name" {
  default = "tf-guide-cosmo-feature-flags"
}

variable "federated_graph_routing_url" {
  default = "http://example.com/graphql"
}

variable "federated_graph_label_matchers" {
  default = ["team=cosmo"]
}

variable "subgraph_name" {
  type    = string
  default = "default_subgraph_name"
}

variable "subgraph_routing_url" {
  type    = string
  default = "http://example.com/subgraph_routing"
}

variable "subgraph_schema" {
  type    = string
  default = "type Query{ b: String }"
}

variable "subgraph_labels" {
  type = map(string)
  default = {
    team = "cosmo"
  }
}

variable "feature_subgraph_name" {
  type    = string
  default = "default_feature_subgraph_name"
}

variable "feature_subgraph_readme" {
  type    = string
  default = "test readme"
}

variable "feature_subgraph_routing_url" {
  type    = string
  default = "http://example.com/feature_subgraph_routing"
}

variable "feature_subgraph_schema" {
  type    = string
  default = "type Query{ c: String }"
}

variable "subscription_protocol" {
  type    = string
  default = "ws"
}

variable "subscription_url" {
  type    = string
  default = "http://example.com/subscription"
}

variable "websocket_subprotocol" {
  type    = string
  default = "graphql-ws"
}

variable "feature_flag_name" {
  type    = string
  default = "default_feature_flag_name"
}

variable "feature_flag_is_enabled" {
  type    = bool
  default = true
}

variable "feature_flag_labels" {
  type = map(string)
  default = {
    team = "cosmo"
  }
}

variable "api_url" {
  type    = string
  default = "http://example.com/graphql"
}

variable "api_key" {
  type = string
}