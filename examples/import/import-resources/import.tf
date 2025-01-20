# file to showcase how to import various resources

# 1. define the resources to import
# 2. find the resource IDs in the UI or wgc
# 3. Run terraform import

# terraform import cosmo_namespace.import <namespace_id>
resource "cosmo_namespace" "import" {}

# terraform import cosmo_contract.import <contract_id>
resource "cosmo_contract" "import" {}

# terraform import cosmo_federated_graph.import <federated_graph_id>
resource "cosmo_federated_graph" "import" {}

# terraform import cosmo_subgraph.import <subgraph_id>
resource "cosmo_subgraph" "import" {}

