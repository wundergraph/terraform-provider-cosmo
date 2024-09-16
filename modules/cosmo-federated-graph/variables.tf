variable "namespace" {
  type        = string
  description = "The name of the namespace to be used for the federated graph"
}

variable "federated_graph" {
  type = object({
    name           = string
    routing_url    = string
    label_matchers = list(string)
  })
  description = "The parameters of the federated graph"
}

variable "subgraphs" {
  type = map(object({
    name        = string
    routing_url = string
    labels      = map(string)
    schema      = string
    readme      = string
  }))
  description = "The subgraphs to be added to the federated graph"
}

variable "router_token_name" {
  type        = string
  description = "The name of the router token to be created"
}
