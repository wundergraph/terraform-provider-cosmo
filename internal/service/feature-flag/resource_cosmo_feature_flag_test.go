package feature_flag_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccFeatureFlags(t *testing.T) {
	fgName, sgName, fsgName :=
		acctest.RandomWithPrefix("test-feature-flag"),
		acctest.RandomWithPrefix("test-subgraph"),
		acctest.RandomWithPrefix("test-feature-subgraph")

	namespace := acctest.RandomWithPrefix("test-namespace")
	ffName := acctest.RandomWithPrefix("test-feature-flag")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureFlagResourceConfig(namespace, fgName, sgName, fsgName, ffName, false),
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
				Config: testAccFeatureFlagResourceConfig(namespace, fgName, sgName, fsgName, ffName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_flag.test", "is_enabled", "true"),
				),
			},
			{
				Config: testAccFeatureFlagResourceConfig(namespace, fgName, sgName, fsgName, "newName", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_flag.test", "name", "newName"),
				),
			},
			{
				ResourceName: "cosmo_feature_flag.test",
				RefreshState: true,
			},
			{
				Config:  testAccFeatureFlagResourceConfig(namespace, fgName, sgName, fsgName, ffName, false),
				Destroy: true,
			},
		},
	})
}

func TestAccFeatureFlagFeatureSubgraphs(t *testing.T) {
	t.Run("Should raise an error when feature_subgraphs is omitted", func(t *testing.T) {
		namespace := acctest.RandomWithPrefix("test-namespace")
		ffName := acctest.RandomWithPrefix("test-feature-flag")

		fgName, sgName, fsgName :=
			acctest.RandomWithPrefix("test-feature-flag"),
			acctest.RandomWithPrefix("test-subgraph"),
			acctest.RandomWithPrefix("test-feature-subgraph")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:      testAccFeatureFlagResourceConfigNoFeatureSubgraphs(namespace, fgName, sgName, fsgName, ffName, false),
					ExpectError: regexp.MustCompile(`.*The argument "feature_subgraphs" is required, but no definition was found.*`),
				},
			},
		})
	})

	t.Run("Should raise an error when feature_subgraphs is empty", func(t *testing.T) {
		namespace := acctest.RandomWithPrefix("test-namespace")
		ffName := acctest.RandomWithPrefix("test-feature-flag")

		fgName, sgName, fsgName :=
			acctest.RandomWithPrefix("test-feature-flag"),
			acctest.RandomWithPrefix("test-subgraph"),
			acctest.RandomWithPrefix("test-feature-subgraph")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
			ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:      testAccFeatureFlagResourceConfigEmptyFeatureSubgraphs(namespace, fgName, sgName, fsgName, ffName, false),
					ExpectError: regexp.MustCompile(`.*Attribute feature_subgraphs set must contain at least 1 elements, got: 0.*`),
				},
			},
		})
	})

}

//nolint:unparam
func testAccFeatureFlagResourceConfig(
	namespace, fgName, sgName, fsgName, ffName string, isEnabled bool) string {
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
`, formatBaseResources(namespace, fgName, sgName, fsgName), ffName, isEnabled)
}

//nolint:unparam
func testAccFeatureFlagResourceConfigNoFeatureSubgraphs(
	namespace, fgName, sgName, fsgName, ffName string, isEnabled bool) string {
	return fmt.Sprintf(`
%s

resource "cosmo_feature_flag" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  labels			  = {
	"team"	= "backend",
	"stage" = "dev"
  }
  is_enabled		  = %t
  depends_on = [cosmo_federated_graph.test, cosmo_feature_subgraph.test]
}
`, formatBaseResources(namespace, fgName, sgName, fsgName), ffName, isEnabled)
}

func testAccFeatureFlagResourceConfigEmptyFeatureSubgraphs(namespace, fgName, sgName, fsgName, ffName string, isEnabled bool) string {
	return fmt.Sprintf(`
%s

resource "cosmo_feature_flag" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  feature_subgraphs   = []
  labels			  = {
	"team"	= "backend",
	"stage" = "dev"
  }
  is_enabled		  = %t
  depends_on = [cosmo_federated_graph.test, cosmo_feature_subgraph.test]
}`, formatBaseResources(namespace, fgName, sgName, fsgName), ffName, isEnabled)

}

func formatBaseResources(namespace, fgName, sgName, fsgName string) string {
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
  readme              =  "test readme"
}

resource "cosmo_feature_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "http://localhost:3000"
  base_subgraph_name  = cosmo_subgraph.test.name
  schema              = <<-EOT
%sEOT
  readme              = "test readme"
}
`, namespace, fgName, sgName,
		acceptance.TestAccValidSubgraphSchema,
		fsgName,
		acceptance.TestAccValidSubgraphSchema,
	)
}
