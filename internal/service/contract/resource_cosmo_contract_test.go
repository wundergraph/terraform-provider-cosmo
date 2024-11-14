package contract_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccContractResource(t *testing.T) {
	name := acctest.RandomWithPrefix("test-contract")
	namespace := acctest.RandomWithPrefix("test-namespace")

	readme := "Initial readme content"

	federatedGraphName := acctest.RandomWithPrefix("test-federated-graph")
	federatedGraphRoutingURL := "https://example.com:3000"

	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	subgraphRoutingURL := "https://subgraph-standalone-example.com"
	subgraphSchema := acceptance.TestAccValidSubgraphSchema

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContractResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, subgraphRoutingURL, subgraphSchema, name, readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_contract.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_contract.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_contract.test", "readme", readme),
				),
			},
			{
				ResourceName: "cosmo_contract.test",
				RefreshState: true,
			},
			{
				Config:  testAccContractResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, subgraphRoutingURL, subgraphSchema, name, readme),
				Destroy: true,
			},
		},
	})
}

func testAccContractResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, subgraphRoutingURL, subgraphSchema, contractName, contractReadme string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "%s"
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
  name      	= "%s"
  namespace 	= cosmo_namespace.test.name
  source     	= cosmo_federated_graph.test.name
  routing_url 	= "http://localhost:3003"
  readme    	= "%s"
}
`, namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, subgraphRoutingURL, subgraphSchema, contractName, contractReadme)
}
