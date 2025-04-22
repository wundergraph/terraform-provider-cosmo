package federated_graph_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccCosmoFederatedGraphImportBasic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-federated-graph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	routingURL := "https://example.com"
	readme := "Initial readme content"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFederatedGraphResourceConfig(namespace, name, routingURL, &readme),
			},
			{
				ResourceName: "cosmo_federated_graph.test",

				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
