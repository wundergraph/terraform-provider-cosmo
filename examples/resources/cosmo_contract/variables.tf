variable "name" {
  type = string
}

variable "routing_url" {
  type = string
}

variable "namespace" {
  type = string
}

variable "source_graph_name" {
  type = string
}

variable "exclude_tags" {
  type = list(string)
  default = []
}

variable "include_tags" {
  type = list(string)
  default = []
}