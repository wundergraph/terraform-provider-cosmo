package federated_graph_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
				Config: testAccFederatedGraphResourceConfig(namespace, name, routingURL, &readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "readme", readme),
				),
			},
			{
				Config: testAccFederatedGraphResourceConfig(namespace, name, routingURL, &newReadme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "readme", newReadme),
				),
			},
			{
				Config: testAccFederatedGraphResourceConfig(namespace, name, updatedRoutingURL, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_federated_graph.test", "routing_url", updatedRoutingURL),
					resource.TestCheckNoResourceAttr("cosmo_federated_graph.test", "readme"),
				),
			},
			{
				ResourceName: "cosmo_federated_graph.test",
				RefreshState: true,
			},
			{
				Config:  testAccFederatedGraphResourceConfig(namespace, name, updatedRoutingURL, nil),
				Destroy: true,
			},
		},
	})
}

func TestAccFederatedGraphResourceInvalidConfig(t *testing.T) {
	name := acctest.RandomWithPrefix("test-federated-graph")
	namespace := acctest.RandomWithPrefix("test-namespace")
	readme := "Initial readme content"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccFederatedGraphResourceConfig(name, namespace, "invalid-url", &readme),
				ExpectError: regexp.MustCompile(`.*Routing URL is not a valid URL*`),
			},
		},
	})
}

func testAccFederatedGraphResourceConfig(namespace, name, routingURL string, readme *string) string {
	var readmePart string
	if readme != nil {
		readmePart = fmt.Sprintf(`readme = "%s"`, *readme)
	}

	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "%s"
  %s
}
`, namespace, name, routingURL, readmePart)
}
