package subgraph_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccSubgraphResource(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")

	federatedGraphName := acctest.RandomWithPrefix("test-subgraph")
	federatedGraphRoutingURL := "https://federated-graph-example.com"

	subgraphName := acctest.RandomWithPrefix("test-subgraph")

	routingURL := "https://subgraph-example.com"
	updatedRoutingURL := "https://updated-subgraph-example.com"

	subgraphSchema := acceptance.TestAccValidSubgraphSchema

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, routingURL, subgraphSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.stage", "dev"),
				),
			},
			{
				Config: testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, updatedRoutingURL, subgraphSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", updatedRoutingURL),
				),
			},
			{
				ResourceName: "cosmo_subgraph.test",
				RefreshState: true,
			},
			{
				Config:  testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, routingURL, subgraphSchema),
				Destroy: true,
			},
		},
	})
}

func TestAccStandaloneSubgraphResource(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")

	federatedGraphName := acctest.RandomWithPrefix("test-subgraph")
	federatedGraphRoutingURL := "https://federated-graph-standalone-subgraph-example.com"

	subgraphName := acctest.RandomWithPrefix("test-subgraph")

	routingURL := "https://subgraph-standalone-example.com"
	subgraphSchema := acceptance.TestAccValidSubgraphSchema

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, routingURL, subgraphSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.stage", "dev"),
				),
			},
			{
				Config: testStandaloneSubgraph(namespace, subgraphName, routingURL, subgraphSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.stage", "dev"),
				),
			},
			{
				ResourceName: "cosmo_subgraph.test",
				RefreshState: true,
			},
			{
				Config:  testStandaloneSubgraph(namespace, subgraphName, routingURL, subgraphSchema),
				Destroy: true,
			},
		},
	})
}

func TestAccSubgraphResourceInvalidSchema(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	subgraphRoutingURL := "https://subgraph-invalid-schema-example.com"

	federatedGraphName := acctest.RandomWithPrefix("test-subgraph")
	federatedGraphRoutingURL := "https://federated-graph-invalid-subgraph-schema-example.com"
	subgraphSchema := "invalid"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, subgraphRoutingURL, subgraphSchema),
				ExpectError: regexp.MustCompile(`.*ERR_INVALID_SUBGRAPH_SCHEMA*`),
			},
		},
	})
}

func TestAccStandaloneSubgraphResourcePublishSchema(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	subgraphRoutingURL := "https://subgraph-publish-schema-example.com"

	subgraphSchema := acceptance.TestAccValidSubgraphSchema
	updatedSubgraphSchema := "invalid"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testStandaloneSubgraph(namespace, subgraphName, subgraphRoutingURL, subgraphSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", subgraphRoutingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.stage", "dev"),
				),
			},
			{
				Config:      testStandaloneSubgraph(namespace, subgraphName, subgraphRoutingURL, updatedSubgraphSchema),
				ExpectError: regexp.MustCompile(`.*ERR_INVALID_SUBGRAPH_SCHEMA*`),
			},
		},
	})
}

func testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphroutingURL, subgraphName, subgraphRoutingURL, subgraphSchema string) string {
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
  %s
  EOT
  labels              = { 
  	"team"	= "backend", 
	"stage" = "dev" 
  }
}
`, namespace, federatedGraphName, federatedGraphroutingURL, subgraphName, subgraphRoutingURL, subgraphSchema)
}

func testStandaloneSubgraph(namespace, subgraphName, subgraphRoutingURL, subgraphSchema string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "%s"
  schema              = <<-EOT
  %s
  EOT
  labels              = { 
  	"team"	= "backend", 
	"stage" = "dev" 
  }
}
`, namespace, subgraphName, subgraphRoutingURL, subgraphSchema)
}
