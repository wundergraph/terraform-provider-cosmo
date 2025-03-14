package monograph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccMonographResource(t *testing.T) {
	name := acctest.RandomWithPrefix("test-monograph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	graphUrl := "http://example.com/graphql"
	rRoutingURL := "http://example.com/routing"
	updatedRoutingURL := "http://example.com/updated-routing"
	updatedGraphURL := "http://example.com/updated-graphql"
	resource.Test(t, resource.TestCase{
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
				Config: testAccMonographResourceConfig(namespace, name, updatedGraphURL, updatedRoutingURL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_monograph.test", "routing_url", updatedRoutingURL),
					resource.TestCheckResourceAttr("cosmo_monograph.test", "graph_url", updatedGraphURL),
				),
			},
			{
				ResourceName: "cosmo_monograph.test",
				RefreshState: true,
			},
			{
				Config:  testAccMonographResourceConfig(namespace, name, graphUrl, rRoutingURL),
				Destroy: true,
			},
		},
	})
}

func testAccMonographResourceConfig(namespace, name, graphURL, routingURL string) string {
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
`, namespace, name, graphURL, routingURL)
}
