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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContractDataSourceConfig(namespace, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cosmo_contract.test", "name", name),
					resource.TestCheckResourceAttr("data.cosmo_contract.test", "namespace", namespace),
				),
			},
			{
				ResourceName: "data.cosmo_contract.test",
				RefreshState: true,
			},
		},
	})
}

func testAccContractDataSourceConfig(namespace, name string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_monograph" "source_graph" {
  name      	= "source-graph"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "https://example.com"
  graph_url 	= "https://example.com"
}

resource "cosmo_contract" "test" {
  name            = "%s"
  namespace       = cosmo_namespace.test.name
  source          = cosmo_monograph.source_graph.name
  routing_url     = "https://example.com"
  readme          = "Initial readme content"
}

data "cosmo_contract" "test" {
  name      = cosmo_contract.test.name
  namespace = cosmo_contract.test.namespace
}
`, namespace, name)
}
