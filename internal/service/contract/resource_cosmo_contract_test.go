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
	federatedGraphroutingURL := "https://example.com:3000"

	graphUrl := "http://example.com/graphql"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContractResourceConfig(namespace, federatedGraphName, federatedGraphroutingURL, graphUrl, name, readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_contract.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_contract.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_contract.test", "readme", readme),
				),
			},
		},
	})
}

func testAccContractResourceConfig(namespace, federatedGraphName, federatedGraphroutingURL, graphUrl, contractName, contractReadme string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_monograph" "test" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "%s"
  graph_url 	= "%s"
}

resource "cosmo_contract" "test" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test.name
  source     	= cosmo_monograph.test.name
  routing_url 	= "http://localhost:3003"
  readme    	= "%s"
}
`, namespace, federatedGraphName, federatedGraphroutingURL, graphUrl, contractName, contractReadme)
}
