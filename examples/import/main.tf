locals {
  schema = <<EOF
type Query {
  getWunder: Wunder!
  getWunderByTitle(title: String!): Wunder!
}
type Wunder {
  title: String
}

type Graph {
  name: String
  wunders: [Wunder]
}
EOF
}

resource "cosmo_namespace" "test" {
  name = "import-test"
}

resource "cosmo_federated_graph" "test" {

  name      = "import-test"
  namespace = cosmo_namespace.test.name

  routing_url    = "http://localhost:9999"
  label_matchers = ["team=backend", "stage=import-test"]
  depends_on     = [cosmo_subgraph.test]
}

// create each stages subgraph
resource "cosmo_subgraph" "test" {
  name   = "import-test-sg"
  schema = local.schema

  namespace = cosmo_namespace.test.name
  labels = {
    "team"  = "backend"
    "stage" = "import-test"
  }

  routing_url = "http://localhost:9997"
}

resource "cosmo_contract" "test" {
  name      = "import-test-contract"
  namespace = cosmo_namespace.test.name
  source    = cosmo_federated_graph.test.name

  routing_url  = "http://localhost:9998"
  exclude_tags = ["backend"]
}