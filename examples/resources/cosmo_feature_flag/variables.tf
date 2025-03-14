variable "name" {
  type = string
}

variable "namespace" {
  type = string
}

variable "feature_subgraphs" {
  type = list(string)
}

variable "labels" {
  type = map(string)
}

variable "is_enabled" {
  type = bool
}
