package feature_subgraph_test

import (
	"fmt"
	"testing"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFeatureSubgraph(t *testing.T) {
	name := acctest.RandomWithPrefix("test-feature-subgraph")
	fgName := acctest.RandomWithPrefix("test-feature-subgraph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	routingURL := "https://example.com"
	updatedRoutingURL := "https://updated-subgraph-example.com"

	subgraphSchema := acceptance.TestAccValidSubgraphSchema
	readme := "Initial readme content"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccFeatureSubgraphResourceConfig(namespace, name, routingURL, name, routingURL, subgraphSchema, readme, fgName, routingURL, subgraphSchema, readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "name", fgName),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "subscription_protocol", api.GraphQLSubscriptionProtocolWS),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "websocket_subprotocol", api.GraphQLWebsocketSubprotocolDefault),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "readme", readme),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "schema", subgraphSchema),
				),
			},
			{
				Config: testAccFeatureSubgraphResourceConfig(namespace, name, routingURL, name, routingURL, subgraphSchema, readme, fgName, routingURL, subgraphSchema, "Updated readme content"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "readme", "Updated readme content"),
				),
			},
			{
				Config: testAccFeatureSubgraphResourceConfig(namespace, name, routingURL, name, routingURL, subgraphSchema, readme, fgName, updatedRoutingURL, subgraphSchema, "Updated readme content"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "routing_url", updatedRoutingURL),
				),
			},
			{
				ResourceName: "cosmo_feature_subgraph.test",
				RefreshState: true,
			},
			{
				Config:  testAccFeatureSubgraphResourceConfig(namespace, name, routingURL, name, routingURL, subgraphSchema, readme, fgName, routingURL, subgraphSchema, readme),
				Destroy: true,
			},
		},
	})
}

func TestAccFeatureSubgraphWithoutSchema(t *testing.T) {
	name := acctest.RandomWithPrefix("test-feature-subgraph")
	fgName := acctest.RandomWithPrefix("test-feature-subgraph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	routingURL := "https://example.com"

	subgraphSchema := acceptance.TestAccValidSubgraphSchema
	readme := "Initial readme content"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccFeatureSubgraphResourceConfig(namespace, name, routingURL, name, routingURL, subgraphSchema, readme, fgName, routingURL, "", readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "name", fgName),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "subscription_protocol", api.GraphQLSubscriptionProtocolWS),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "websocket_subprotocol", api.GraphQLWebsocketSubprotocolDefault),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "readme", readme),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "schema", ""),
				),
			},
			{
				Config: testAccFeatureSubgraphResourceConfig(namespace, name, routingURL, name, routingURL, subgraphSchema, readme, fgName, routingURL, subgraphSchema, "Updated readme content"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "readme", "Updated readme content"),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "schema", subgraphSchema),
				),
			},
			{
				ResourceName: "cosmo_feature_subgraph.test",
				RefreshState: true,
			},
			{
				Config:  testAccFeatureSubgraphResourceConfig(namespace, name, routingURL, name, routingURL, subgraphSchema, readme, fgName, routingURL, subgraphSchema, readme),
				Destroy: true,
			},
		},
	})
}

func testAccFeatureSubgraphResourceConfig(
	namespace, federatedGraphName, federatedGraphroutingURL,
	subgraphName, subgraphRoutingURL, subgraphSchema, readme,
	fgName, fgRoutingURL, fgSchema, fgReadme string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "%s"
  label_matchers = ["team=backend"]

  depends_on = [cosmo_subgraph.test]
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "%s"
  schema              = <<-EOT
%sEOT
  labels              = { 
  	"team"	= "backend", 
	"stage" = "dev" 
  }
  readme              =  "%s"
}

resource "cosmo_feature_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "%s"
  base_subgraph_name  = cosmo_subgraph.test.name
  schema              = <<-EOT
%sEOT
  readme              = "%s"
}
`, namespace, federatedGraphName, federatedGraphroutingURL, subgraphName, subgraphRoutingURL, subgraphSchema, readme, fgName, fgRoutingURL, fgSchema, fgReadme)
}
