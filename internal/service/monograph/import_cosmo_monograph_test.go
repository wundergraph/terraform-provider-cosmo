package monograph_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
	"testing"
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
				ResourceName:            "cosmo_monograph.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"graph_url"}, // TODO we currently don't get this value from the API
			},
		},
	})
}
