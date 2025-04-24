package subgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccSubgraphDataSource(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")

	federatedGraphName := acctest.RandomWithPrefix("test-federated-graph")
	federatedGraphRoutingURL := "https://example.com"

	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	subgraphRoutingURL := "https://example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubgraphDataSourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, subgraphRoutingURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("data.cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("data.cosmo_subgraph.test", "routing_url", subgraphRoutingURL),
					resource.TestCheckResourceAttr("data.cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("data.cosmo_subgraph.test", "labels.stage", "dev"),
					resource.TestCheckNoResourceAttr("data.cosmo_subgraph.test", "subscription_url"),
					resource.TestCheckResourceAttr("data.cosmo_subgraph.test", "subscription_protocol", "ws"),
					resource.TestCheckResourceAttr("data.cosmo_subgraph.test", "websocket_subprotocol", "auto"),
					resource.TestCheckResourceAttr("data.cosmo_subgraph.test", "is_event_driven_graph", "false"),
				),
			},
		},
	})
}

func testAccSubgraphDataSourceConfig(namespace, federatedGraphName, federatedGraphroutingURL, subgraphName, subgraphRoutingURL string) string {
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
  labels              = { "team" = "backend", "stage" = "dev" }
}

data "cosmo_subgraph" "test" {
  name      = cosmo_subgraph.test.name
  namespace = cosmo_subgraph.test.namespace
}
`, namespace, federatedGraphName, federatedGraphroutingURL, subgraphName, subgraphRoutingURL)
}
