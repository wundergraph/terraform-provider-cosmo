package contract_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccContractDataSource(t *testing.T) {
	name := acctest.RandomWithPrefix("test-contract")
	namespace := acctest.RandomWithPrefix("test-namespace")

	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	subgraphRoutingURL := "https://subgraph-standalone-example.com"
	subgraphSchema := acceptance.TestAccValidSubgraphSchema

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContractDataSourceConfig(namespace, subgraphName, subgraphRoutingURL, subgraphSchema, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cosmo_contract.test", "name", name),
					resource.TestCheckResourceAttr("data.cosmo_contract.test", "namespace", namespace),
				),
			},
			{
				ResourceName: "data.cosmo_contract.test",
				RefreshState: true,
			},
			{
				Config:  testAccContractDataSourceConfig(namespace, subgraphName, subgraphRoutingURL, subgraphSchema, name),
				Destroy: true,
			},
		},
	})
}

func testAccContractDataSourceConfig(namespace, subgraphName, subgraphRoutingURL, subgraphSchema, name string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "source_graph" {
  name      	= "source-graph"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "https://example.com"
  depends_on = [cosmo_subgraph.test]
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "%s"
  schema              = <<-EOT
  %s
  EOT
  labels = {}
}

resource "cosmo_contract" "test" {
  name            = "%s"
  namespace       = cosmo_namespace.test.name
  source          = cosmo_federated_graph.source_graph.name
  routing_url     = "https://example.com"
  readme          = "Initial readme content"
}

data "cosmo_contract" "test" {
  name      = cosmo_contract.test.name
  namespace = cosmo_contract.test.namespace
}
`, namespace, subgraphName, subgraphRoutingURL, subgraphSchema, name)
}
