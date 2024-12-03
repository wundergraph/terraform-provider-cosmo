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

	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	subgraphRoutingURL := "https://subgraph-standalone-example.com"
	subgraphSchema := acceptance.TestAccValidSubgraphSchema

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContractResourceConfig(namespace, federatedGraphName, subgraphName, subgraphRoutingURL, subgraphSchema, name, "internal", &readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_contract.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_contract.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_contract.test", "readme", readme),
					resource.TestCheckResourceAttr("cosmo_contract.test", "readme", readme),
					resource.TestCheckResourceAttr("cosmo_contract.test", "exclude_tags.0", "internal"),
				),
			},
			{
				Config: testAccContractResourceConfig(namespace, federatedGraphName, subgraphName, subgraphRoutingURL, subgraphSchema, name, "external", &readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_contract.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_contract.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_contract.test", "readme", readme),
					resource.TestCheckResourceAttr("cosmo_contract.test", "exclude_tags.0", "external"),
				),
			},
			{
				ResourceName: "cosmo_contract.test",
				RefreshState: true,
			},
			{
				Config:  testAccContractResourceConfig(namespace, federatedGraphName, subgraphName, subgraphRoutingURL, subgraphSchema, name, "external", &readme),
				Destroy: true,
			},
			{
				Config: testAccContractOfMonographResourceConfig(namespace, federatedGraphName, subgraphRoutingURL, subgraphSchema, name, readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_contract.test_mono", "name", name),
					resource.TestCheckResourceAttr("cosmo_contract.test_mono", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_contract.test_mono", "readme", readme),
				),
			},
			{
				Config:  testAccContractOfMonographResourceConfig(namespace, federatedGraphName, subgraphRoutingURL, subgraphSchema, name, readme),
				Destroy: true,
			},
		},
	})
}

func TestOptionalValuesOfContractResource(t *testing.T) {
	name := acctest.RandomWithPrefix("test-contract")
	namespace := acctest.RandomWithPrefix("test-namespace")

	federatedGraphName := acctest.RandomWithPrefix("test-federated-graph")

	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	subgraphRoutingURL := "https://subgraph-standalone-example.com"
	subgraphSchema := acceptance.TestAccValidSubgraphSchema

	readme := "Initial readme content"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContractResourceConfig(namespace, federatedGraphName, subgraphName, subgraphRoutingURL, subgraphSchema, name, "internal", &readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_contract.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_contract.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_contract.test", "exclude_tags.0", "internal"),
					resource.TestCheckResourceAttr("cosmo_contract.test", "readme", readme),
				),
			},
			{
				Config: testAccContractResourceConfig(namespace, federatedGraphName, subgraphName, subgraphRoutingURL, subgraphSchema, name, "external", nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_contract.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_contract.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_contract.test", "exclude_tags.0", "external"),
					resource.TestCheckNoResourceAttr("cosmo_contract.test", "readme"),
				),
			},
		},
	})
}

func testAccContractResourceConfig(namespace, federatedGraphName, subgraphName, subgraphRoutingURL, subgraphSchema, contractName, contractExcludeTag string, contractReadme *string) string {
	var readmePart string
	if contractReadme != nil {
		readmePart = fmt.Sprintf(`readme = "%s"`, *contractReadme)
	}

	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "https://example.com:3000"
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
  exclude_tags = ["%s"]
  %s
}
`, namespace, federatedGraphName, subgraphName, subgraphRoutingURL, subgraphSchema, contractName, contractExcludeTag, readmePart)
}

func testAccContractOfMonographResourceConfig(namespace, federatedGraphName, graphUrl, schema, contractName, contractReadme string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test_mono" {
  name = "%s"
}

resource "cosmo_monograph" "test_mono" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test_mono.name
  routing_url 	= "https://example.com:3000"
  graph_url 	= "%s"
  schema              = <<-EOT
  %s
  EOT
}

resource "cosmo_contract" "test_mono" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test_mono.name
  source     	= cosmo_monograph.test_mono.name
  routing_url 	= "http://localhost:3003"
  readme    	= "%s"
}
`, namespace, federatedGraphName, graphUrl, schema, contractName, contractReadme)
}
