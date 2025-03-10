package feature_subgraph_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
)

func TestAccImportFeatureSubgraph(t *testing.T) {
	fsgName := acctest.RandomWithPrefix("test-feature-subgraph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	fgName, sgName :=
		acctest.RandomWithPrefix("test-federated-graph"),
		acctest.RandomWithPrefix("test-subgraph")

	routingURL := "https://example.com"

	subgraphSchema := acceptance.TestAccValidSubgraphSchema
	readme := "Initial readme content"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccFeatureSubgraphResourceConfig(namespace, fgName, sgName, fsgName, routingURL, subgraphSchema, readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "name", fsgName),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "subscription_protocol", api.GraphQLSubscriptionProtocolWS),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "websocket_subprotocol", api.GraphQLWebsocketSubprotocolDefault),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "readme", readme),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "schema", subgraphSchema),
				),
			},
			{
				// Import the resource by ID
				ResourceName:      "cosmo_feature_subgraph.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccFeatureSubgraphResourceConfig(namespace, fgName, sgName, fsgName, routingURL, subgraphSchema, readme),
				Destroy: true,
			},
		},
	})
}
