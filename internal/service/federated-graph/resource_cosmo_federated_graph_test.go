package federated_graph_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccFederatedGraphResource(t *testing.T) {
	name := acctest.RandomWithPrefix("test-federated-graph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	routingURL := "https://example.com"
	updatedRoutingURL := "https://updated-example.com"

	readme := "Initial readme content"
	newReadme := "Updated readme content"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFederatedGraphResourceConfig(namespace, name, routingURL, readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "readme", readme),
				),
			},
			{
				Config: testAccFederatedGraphResourceConfig(namespace, name, routingURL, newReadme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "readme", newReadme),
				),
			},
			{
				Config: testAccFederatedGraphResourceConfig(namespace, name, updatedRoutingURL, newReadme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "routing_url", updatedRoutingURL),
				),
			},
		},
	})
}

func TestAccFederatedGraphResourceInvalidConfig(t *testing.T) {
	name := acctest.RandomWithPrefix("test-federated-graph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccFederatedGraphResourceConfig(name, namespace, "invalid-url", ""),
				ExpectError: regexp.MustCompile(`.*failed to create resource*`),
			},
		},
	})
}

func testAccFederatedGraphResourceConfig(namespace, name, routingURL, readme string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "%s"
  readme    	= "%s"
}
`, namespace, name, routingURL, readme)
}
