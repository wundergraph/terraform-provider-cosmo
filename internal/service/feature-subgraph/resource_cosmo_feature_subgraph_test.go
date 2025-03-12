package feature_subgraph_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
)

func TestAccFeatureSubgraph(t *testing.T) {
	fsgName := acctest.RandomWithPrefix("test-feature-subgraph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	fgName, sgName :=
		acctest.RandomWithPrefix("test-federated-graph"),
		acctest.RandomWithPrefix("test-subgraph")

	routingURL := "https://example.com"
	updatedRoutingURL := "https://updated-subgraph-example.com"

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
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "base_subgraph_name", sgName),
				),
			},
			{
				Config: testAccFeatureSubgraphResourceConfig(namespace, fgName, sgName, fsgName, routingURL, subgraphSchema, "Updated readme content"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "readme", "Updated readme content"),
				),
			},
			{
				Config: testAccFeatureSubgraphResourceConfig(namespace, fgName, sgName, fsgName, updatedRoutingURL, subgraphSchema, "Updated readme content"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "routing_url", updatedRoutingURL),
				),
			},
			{
				ResourceName: "cosmo_feature_subgraph.test",
				RefreshState: true,
			},
			{
				Config:  testAccFeatureSubgraphResourceConfig(namespace, fgName, sgName, fsgName, routingURL, subgraphSchema, readme),
				Destroy: true,
			},
		},
	})
}

func TestAccFeatureSubgraphWithoutSchema(t *testing.T) {
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
				Config: testAccFeatureSubgraphResourceConfig(namespace, fgName, sgName, fsgName, routingURL, "", readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "name", fsgName),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "subscription_protocol", api.GraphQLSubscriptionProtocolWS),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "websocket_subprotocol", api.GraphQLWebsocketSubprotocolDefault),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "readme", readme),
					resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "schema", ""),
				),
			},
			{
				Config: testAccFeatureSubgraphResourceConfig(namespace, fgName, sgName, fsgName, routingURL, subgraphSchema, "Updated readme content"),
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
				Config:  testAccFeatureSubgraphResourceConfig(namespace, fgName, sgName, fsgName, routingURL, subgraphSchema, readme),
				Destroy: true,
			},
		},
	})
}

func TestAccFeatureSubgraphWithoutBaseSubgraphName(t *testing.T) {
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
				Config:      testAccFeatureSubgraphResourceConfigNoBaseSubgraph(namespace, fgName, sgName, fsgName, routingURL, subgraphSchema, readme),
				ExpectError: regexp.MustCompile(".*A feature subgraph requires a base.*"),
			},
		},
	})
}

// TODO: uncomment when URL normalization was fixed in the control plane
//func TestAccFeatureSubgraphRoutingURL(t *testing.T) {
//	t.Run("should contain the correct state when providing localhost without protocol", func(t *testing.T) {
//		fsgName := acctest.RandomWithPrefix("test-feature-subgraph")
//		namespace := acctest.RandomWithPrefix("test-namespace")
//
//		fgName, sgName :=
//			acctest.RandomWithPrefix("test-federated-graph"),
//			acctest.RandomWithPrefix("test-subgraph")
//
//		readme := "Initial readme content"
//
//		routingURL := "localhost:3000"
//
//		resource.ParallelTest(t, resource.TestCase{
//			PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
//			ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
//
//			Steps: []resource.TestStep{
//				{
//					Config: testAccFeatureSubgraphResourceConfig(namespace, fgName, sgName, fsgName, routingURL, "", readme),
//					Check: resource.ComposeTestCheckFunc(
//						resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "name", fsgName),
//						resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "namespace", namespace),
//						resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "routing_url", "localhost:3000"),
//						resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "subscription_protocol", api.GraphQLSubscriptionProtocolWS),
//						resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "websocket_subprotocol", api.GraphQLWebsocketSubprotocolDefault),
//						resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "readme", readme),
//						resource.TestCheckResourceAttr("cosmo_feature_subgraph.test", "schema", ""),
//					),
//				},
//				{
//					Config:  testAccFeatureSubgraphResourceConfig(namespace, fgName, sgName, fsgName, routingURL, "", readme),
//					Destroy: true,
//				},
//			},
//		})
//	})
//}

func formatBaseResource(namespace, fgName, sgName string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "http://localhost:3000"
  label_matchers = ["team=backend"]

  depends_on = [cosmo_subgraph.test]
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "http://localhost:3000"
  schema              = <<-EOT
%sEOT
  labels              = { 
  	"team"	= "backend", 
	"stage" = "dev" 
  }
  readme              =  "Test Readme"
}
`, namespace, fgName, sgName, acceptance.TestAccValidSubgraphSchema)
}

//nolint:unparam
func testAccFeatureSubgraphResourceConfig(
	namespace, fgName, sgName, fsgName, fsgRoutingURL, fsgSchema, fsgReadme string) string {
	return fmt.Sprintf(`
%s

resource "cosmo_feature_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "%s"
  base_subgraph_name  = cosmo_subgraph.test.name
  schema              = <<-EOT
%sEOT
  readme              = "%s"

  depends_on = [cosmo_subgraph.test]
}
`, formatBaseResource(namespace, fgName, sgName), fsgName, fsgRoutingURL, fsgSchema, fsgReadme)
}

func testAccFeatureSubgraphResourceConfigNoBaseSubgraph(
	namespace, fgName, sgName, fsgName, fsgRoutingURL, fsgSchema, fsgReadme string) string {
	return fmt.Sprintf(`
%s

resource "cosmo_feature_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "%s"
  base_subgraph_name  = ""
  schema              = <<-EOT
%sEOT
  readme              = "%s"

  depends_on = [cosmo_subgraph.test]
}
`, formatBaseResource(namespace, fgName, sgName), fsgName, fsgRoutingURL, fsgSchema, fsgReadme)
}
