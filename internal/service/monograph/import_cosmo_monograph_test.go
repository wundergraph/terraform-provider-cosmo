package monograph_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccCosmoMonographGraphImportBasic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-monograph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	graphUrl := "http://example.com/graphql"
	rRoutingURL := "http://example.com/routing"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMonographResourceConfig(namespace, name, graphUrl, rRoutingURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_monograph.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_monograph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_monograph.test", "graph_url", graphUrl),
					resource.TestCheckResourceAttr("cosmo_monograph.test", "routing_url", rRoutingURL),
					resource.TestCheckResourceAttrSet("cosmo_monograph.test", "id"),
				),
			},
			{
				ResourceName:      "cosmo_monograph.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccImportCosmoMonographByProvidingFederatedGraphId(t *testing.T) {
	name := acctest.RandomWithPrefix("test-monograph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	graphUrl := "http://example.com/graphql"
	rRoutingURL := "http://example.com/routing"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFederatedGraphResourceConfig(namespace, name, graphUrl, rRoutingURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_monograph.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_monograph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_monograph.test", "graph_url", graphUrl),
					resource.TestCheckResourceAttr("cosmo_monograph.test", "routing_url", rRoutingURL),
					resource.TestCheckResourceAttrSet("cosmo_monograph.test", "id"),
				),
			},
			{
				ResourceName:      "cosmo_monograph.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName: "cosmo_monograph.test",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					federatedGraph, ok := s.RootModule().Resources["cosmo_federated_graph.test"]
					if !ok {
						return "", errors.New("federated graph not found")
					}

					return federatedGraph.Primary.ID, nil
				},
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile(`.*A monograph must be a non-federated.*`),
			},
		},
	})
}

func testAccFederatedGraphResourceConfig(namespace, name, graphUrl, routingURL string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_monograph" "test" {
	name       = "%s"
	namespace  = cosmo_namespace.test.name
	graph_url  = "%s"
	routing_url = "%s"
}

resource "cosmo_federated_graph" "test" {
	name      = "%s-federated-graph"
	namespace = cosmo_namespace.test.name
	routing_url = "%s"
}
`, namespace, name, graphUrl, routingURL, name, routingURL)
}
