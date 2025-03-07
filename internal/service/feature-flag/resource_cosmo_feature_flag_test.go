package feature_flag_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccFeatureFlags(t *testing.T) {
	name := acctest.RandomWithPrefix("test-feature-subgraph")
	fgName := acctest.RandomWithPrefix("test-feature-subgraph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	routingURL := "https://example.com"

	subgraphSchema := acceptance.TestAccValidSubgraphSchema
	readme := "Initial readme content"

	ffName := acctest.RandomWithPrefix("test-feature-flag")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureFlagResourceConfig(namespace, name, routingURL, name, routingURL, subgraphSchema, readme, fgName, routingURL, subgraphSchema, readme, ffName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_flag.test", "name", ffName),
					resource.TestCheckResourceAttr("cosmo_feature_flag.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_feature_flag.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_feature_flag.test", "labels.stage", "dev"),
					resource.TestCheckResourceAttr("cosmo_feature_flag.test", "is_enabled", "false"),
					resource.TestCheckResourceAttrSet("cosmo_feature_flag.test", "feature_subgraphs.0"),
				),
			},
			{
				Config: testAccFeatureFlagResourceConfig(namespace, name, routingURL, name, routingURL, subgraphSchema, readme, fgName, routingURL, subgraphSchema, readme, ffName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_flag.test", "is_enabled", "true"),
				),
			},
			{
				Config: testAccFeatureFlagResourceConfig(namespace, name, routingURL, name, routingURL, subgraphSchema, readme, fgName, routingURL, subgraphSchema, readme, "newName", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_flag.test", "name", "newName"),
				),
			},
			{
				ResourceName: "cosmo_feature_flag.test",
				RefreshState: true,
			},
			{
				Config:  testAccFeatureFlagResourceConfig(namespace, name, routingURL, name, routingURL, subgraphSchema, readme, fgName, routingURL, subgraphSchema, readme, ffName, false),
				Destroy: true,
			},
		},
	})

}

func testAccFeatureFlagResourceConfig(
	namespace, federatedGraphName, federatedGraphroutingURL,
	subgraphName, subgraphRoutingURL, subgraphSchema, readme,
	fgName, fgRoutingURL, fgSchema, fgReadme,
	ffName string, isEnabled bool) string {
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

resource "cosmo_feature_flag" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  feature_subgraphs   = [cosmo_feature_subgraph.test.name]
  labels			  = {
	"team"	= "backend",
	"stage" = "dev"
  }
  is_enabled		  = %t
  depends_on = [cosmo_federated_graph.test, cosmo_feature_subgraph.test]
}
`, namespace, federatedGraphName, federatedGraphroutingURL, subgraphName, subgraphRoutingURL,
		subgraphSchema, readme, fgName, fgRoutingURL, fgSchema, fgReadme,
		ffName, isEnabled)
}
