package clients_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccClientsDataSource(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")
	graphName := acctest.RandomWithPrefix("test-graph")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	clientName := acctest.RandomWithPrefix("test-client")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClientsDataSourceConfig(namespace, graphName, subgraphName, clientName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cosmo_clients.test", "federated_graph_name", graphName),
					resource.TestCheckResourceAttr("data.cosmo_clients.test", "clients.#", "1"),
					resource.TestCheckResourceAttr("data.cosmo_clients.test", "clients.0.name", clientName),
					resource.TestCheckResourceAttrSet("data.cosmo_clients.test", "clients.0.id"),
				),
			},
		},
	})
}

func testAccClientsDataSourceConfig(namespace, graphName, subgraphName, clientName string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name           = "%s"
  namespace      = cosmo_namespace.test.name
  routing_url    = "https://example.com"
  label_matchers = ["team=backend"]

  depends_on = [cosmo_subgraph.test]
}

resource "cosmo_subgraph" "test" {
  name        = "%s"
  namespace   = cosmo_namespace.test.name
  routing_url = "https://subgraph-example.com"
  schema      = <<-EOT
  %s
  EOT
  labels = {
    "team" = "backend"
  }
}

resource "cosmo_persisted_operations" "test" {
  federated_graph_name = cosmo_federated_graph.test.name
  namespace            = cosmo_namespace.test.name
  client_name          = "%s"
  operations = {
    capsules = "query Capsules { capsules { id } }"
  }
}

data "cosmo_clients" "test" {
  federated_graph_name = cosmo_federated_graph.test.name
  namespace            = cosmo_namespace.test.name

  depends_on = [cosmo_persisted_operations.test]
}
`, namespace, graphName, subgraphName, acceptance.TestAccValidSubgraphSchema, clientName)
}
