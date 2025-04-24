package feature_flag_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccFeatureFlagDataSource(t *testing.T) {
	t.Run("Should import datasource and refresh the state without differences", func(t *testing.T) {

		namespace := acctest.RandomWithPrefix("test-namespace")
		fgName := acctest.RandomWithPrefix("test-feature-flag")
		sgName := acctest.RandomWithPrefix("test-subgraph")
		fsgName := acctest.RandomWithPrefix("test-feature-subgraph")
		ffName := acctest.RandomWithPrefix("test-feature-flag")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,

			Steps: []resource.TestStep{
				{
					Config: testAccFeatureFlagDataSourceConfig(namespace, fgName, sgName, fsgName, ffName, false),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.cosmo_feature_flag.test", "name", ffName),
						resource.TestCheckResourceAttr("data.cosmo_feature_flag.test", "namespace", namespace),
						resource.TestCheckResourceAttr("data.cosmo_feature_flag.test", "labels.team", "backend"),
						resource.TestCheckResourceAttr("data.cosmo_feature_flag.test", "labels.stage", "dev"),
						resource.TestCheckResourceAttr("data.cosmo_feature_flag.test", "is_enabled", "false"),
						resource.TestCheckResourceAttrSet("data.cosmo_feature_flag.test", "feature_subgraphs.0"),
					),
				},
				{
					ResourceName: "data.cosmo_feature_flag.test",
					RefreshState: true,
				},
				{
					Config:  testAccFeatureFlagDataSourceConfig(namespace, fgName, sgName, fsgName, ffName, true),
					Destroy: true,
				},
			},
		})
	})

	t.Run("Should import datasource without namespace and use default instead", func(t *testing.T) {
		fgName := acctest.RandomWithPrefix("test-feature-flag")
		sgName := acctest.RandomWithPrefix("test-subgraph")
		fsgName := acctest.RandomWithPrefix("test-feature-subgraph")
		ffName := acctest.RandomWithPrefix("test-feature-flag")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,

			Steps: []resource.TestStep{
				{
					Config: testAccFeatureFlagDataSourceConfigWithDefaultNamespace(fgName, sgName, fsgName, ffName, false),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.cosmo_feature_flag.test", "name", ffName),
						resource.TestCheckResourceAttr("data.cosmo_feature_flag.test", "namespace", "default"),
						resource.TestCheckResourceAttr("data.cosmo_feature_flag.test", "labels.team", "backend"),
						resource.TestCheckResourceAttr("data.cosmo_feature_flag.test", "labels.stage", "dev"),
						resource.TestCheckResourceAttr("data.cosmo_feature_flag.test", "is_enabled", "false"),
						resource.TestCheckResourceAttrSet("data.cosmo_feature_flag.test", "feature_subgraphs.0"),
					),
				},
				{
					ResourceName: "data.cosmo_feature_flag.test",
					RefreshState: true,
				},
				{
					Config:  testAccFeatureFlagDataSourceConfigWithDefaultNamespace(fgName, sgName, fsgName, ffName, false),
					Destroy: true,
				},
			},
		})
	})
}

func testAccFeatureFlagDataSourceConfigWithDefaultNamespace(fgName, sgName, fsgName, ffName string, isEnabled bool) string {
	return fmt.Sprintf(`
resource "cosmo_federated_graph" "test" {
  name      	= "%s"
  namespace 	= "default"
  routing_url 	= "https://example.com"
  depends_on = [cosmo_subgraph.test]
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  namespace           = "default"
  routing_url         = "https://subgraph-standalone-example.com"
  schema              = <<-EOT
%sEOT
  labels = {
	"team" = "backend"
    "stage" = "dev"
  }
}

resource "cosmo_feature_subgraph" "test" {
  name                = "%s"
  namespace           = "default"
  routing_url         = "http://localhost:3000"
  base_subgraph_name  = cosmo_subgraph.test.name
  schema              = <<-EOT
%sEOT
  readme              = "test readme"
}

resource "cosmo_feature_flag" "test" {
  name                = "%s"
  namespace           = "default"
  feature_subgraphs   = [cosmo_feature_subgraph.test.name]
  labels			  = {
    "team"	= "backend",
    "stage" = "dev"
  }
  is_enabled		  = %t
  depends_on = [cosmo_federated_graph.test, cosmo_subgraph.test]
}

data "cosmo_feature_flag" "test" {
  name = cosmo_feature_flag.test.name
}
`, fgName, sgName, acceptance.TestAccValidSubgraphSchema, fsgName, acceptance.TestAccValidSubgraphSchema, ffName, isEnabled)
}

func testAccFeatureFlagDataSourceConfig(namespace, fgName, sgName, fsgName, ffName string, isEnabled bool) string {
	return fmt.Sprintf(`
%s

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

data "cosmo_feature_flag" "test" {
  name = cosmo_feature_flag.test.name
  namespace = cosmo_feature_flag.test.namespace
}
`, formatBaseResources(namespace, fgName, sgName, fsgName), ffName, isEnabled)

}
