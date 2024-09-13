package subgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccSubgraphResource(t *testing.T) {
	federatedGraphName := acctest.RandomWithPrefix("test-subgraph")
	namespace := acctest.RandomWithPrefix("test-namespace")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	federatedGraphRoutingURL := "https://example.com"

	routingURL := "https://example.com"
	updatedRoutingURL := "https://updated-example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, routingURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.stage", "dev"),
				),
			},
			{
				Config: testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, updatedRoutingURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", updatedRoutingURL),
				),
			},
			{
				ResourceName: "cosmo_subgraph.test",
				RefreshState: true,
			},
		},
	})
}

func testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphroutingURL, subgraphName, subgraphRoutingURL string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "%s"
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "%s"
  labels              = { 
  	"team"	= "backend", 
	"stage" = "dev" 
  }
}
`, namespace, federatedGraphName, federatedGraphroutingURL, subgraphName, subgraphRoutingURL)
}
