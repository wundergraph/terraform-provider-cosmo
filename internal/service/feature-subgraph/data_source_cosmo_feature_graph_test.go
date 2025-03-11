package feature_subgraph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccFeatureSubgraphDataSource(t *testing.T) {
	t.Run("Should create and read a feature subgraph data source", func(t *testing.T) {
		fgName := acctest.RandomWithPrefix("test-federated-graph")
		sgName := acctest.RandomWithPrefix("test-subgraph")
		fsgName := acctest.RandomWithPrefix("test-feature-subgraph")
		namespace := acctest.RandomWithPrefix("test-namespace")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccFeatureSubgraphDataSourceConfig(namespace, fgName, sgName, fsgName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "name", fsgName),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "namespace", namespace),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "routing_url", "https://example.com"),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "base_subgraph_name", sgName),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "subscription_url", "https://subscription.example.com"),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "subscription_protocol", "ws"),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "websocket_subprotocol", "auto"),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "readme", "Initial readme content"),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "schema", acceptance.TestAccValidSubgraphSchema),
					),
				},
				{
					ResourceName: "data.cosmo_feature_subgraph.test",
					RefreshState: true,
				},
				{
					Config:  testAccFeatureSubgraphDataSourceConfig(namespace, fgName, sgName, fsgName),
					Destroy: true,
				},
			},
		})
	})

	t.Run("Should import feature subgraph without namespace", func(t *testing.T) {
		fgName := acctest.RandomWithPrefix("test-federated-graph")
		sgName := acctest.RandomWithPrefix("test-subgraph")
		fsgName := acctest.RandomWithPrefix("test-feature-subgraph")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccFeatureSubgraphDataSourceConfigNoNamespace(fgName, sgName, fsgName),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "name", fsgName),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "namespace", "default"),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "routing_url", "https://example.com"),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "base_subgraph_name", sgName),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "subscription_url", "https://subscription.example.com"),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "subscription_protocol", "ws"),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "websocket_subprotocol", "auto"),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "readme", "Initial readme content"),
						resource.TestCheckResourceAttr("data.cosmo_feature_subgraph.test", "schema", acceptance.TestAccValidSubgraphSchema),
					),
				},
				{
					ResourceName: "data.cosmo_feature_subgraph.test",
					RefreshState: true,
				},
				{
					Config:  testAccFeatureSubgraphDataSourceConfigNoNamespace(fgName, sgName, fsgName),
					Destroy: true,
				},
			},
		})
	})

}

func testAccFeatureSubgraphDataSourceConfigNoNamespace(fgName, sgName, fsgName string) string {
	return fmt.Sprintf(`
resource "cosmo_federated_graph" "test" {
  name      	= "%s"
  namespace 	= "default"
  routing_url 	= "http://localhost:3000"
  label_matchers = ["team=backend"]

  depends_on = [cosmo_subgraph.test]
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  namespace           = "default"
  routing_url         = "http://localhost:3000"
  schema              = <<-EOT
%sEOT
  labels              = { 
  	"team"	= "backend", 
	"stage" = "dev" 
  }
  readme              =  "Test Readme"
}

resource "cosmo_feature_subgraph" "test" {
  name                = "%s"
  namespace           = "default"
  routing_url         = "https://example.com"
  base_subgraph_name  = cosmo_subgraph.test.name
  subscription_url    = "https://subscription.example.com"
  readme              = "Initial readme content"
  schema              = <<-EOT
%sEOT
}

data "cosmo_feature_subgraph" "test" {
  name      = cosmo_feature_subgraph.test.name
  namespace = cosmo_feature_subgraph.test.namespace
}
`, fgName, sgName, acceptance.TestAccValidSubgraphSchema, fsgName, acceptance.TestAccValidSubgraphSchema)
}

func testAccFeatureSubgraphDataSourceConfig(namespace, fgName, sgName, fsgName string) string {
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

resource "cosmo_feature_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "https://example.com"
  base_subgraph_name  = cosmo_subgraph.test.name
  subscription_url    = "https://subscription.example.com"
  readme              = "Initial readme content"
  schema              = <<-EOT
%sEOT
}

data "cosmo_feature_subgraph" "test" {
  name      = cosmo_feature_subgraph.test.name
  namespace = cosmo_feature_subgraph.test.namespace
}
`, namespace, fgName, sgName, acceptance.TestAccValidSubgraphSchema, fsgName, acceptance.TestAccValidSubgraphSchema)
}
