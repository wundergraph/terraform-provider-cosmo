package feature_flags_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccFeatureFlagResource(t *testing.T) {
	name := acctest.RandomWithPrefix("test-feature-flag")
	namespace := acctest.RandomWithPrefix("test-namespace")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	subgraphRoutingURL := "https://example.com"
	subgraphSchema := `type Query { hello: String }`
	isEnabled := true

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureFlagResourceConfig(namespace, name, subgraphName, subgraphRoutingURL, subgraphSchema, isEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_feature_flag.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_feature_flag.test", "namespace", namespace),
				),
			},
		},
	})
}

func testAccFeatureFlagResourceConfig(namespace, name, subgraphName, subgraphRoutingURL, subgraphSchema string, isEnabled bool) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_subgraph" "base" {
  name                = "base"
  routing_url         = "https://example.com"
  namespace           = cosmo_namespace.test.name
  labels              = { 
  	"team"	= "backend", 
  }
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  routing_url         = "%s"
  namespace           = cosmo_namespace.test.name
  base_subgraph_name  = cosmo_subgraph.base.name
  is_feature_subgraph = true
  schema              = <<-EOT
  %s
  EOT
  labels              = { 
  	"team"	= "backend", 
  }
}

resource "cosmo_feature_flag" "test" {
  name                   = "%s"
  is_enabled             = %t
  labels                 = []
  feature_subgraph_names = [cosmo_subgraph.test.name]
  namespace              = cosmo_namespace.test.name
}
`, namespace, subgraphName, subgraphRoutingURL, subgraphSchema, name, isEnabled)
}
