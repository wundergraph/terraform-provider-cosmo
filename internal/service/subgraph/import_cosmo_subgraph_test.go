package subgraph_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
	"testing"
)

func TestAccImportCosmoSubgraphBasic(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")

	federatedGraphName := acctest.RandomWithPrefix("test-subgraph")
	federatedGraphRoutingURL := "https://federated-graph-example.com"

	subgraphName := acctest.RandomWithPrefix("test-subgraph")

	routingURL := "https://subgraph-example.com"

	subgraphSchema := acceptance.TestAccValidSubgraphSchema
	readme := "Initial readme content"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, routingURL, subgraphSchema, readme),
			},
			{
				ResourceName:      "cosmo_subgraph.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
