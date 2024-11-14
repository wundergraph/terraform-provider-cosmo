package federated_graph_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccFederatedGraphDataSource(t *testing.T) {
	name := acctest.RandomWithPrefix("test-federated-graph")
	namespace := acctest.RandomWithPrefix("test-namespace")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFederatedGraphDataSourceConfig(namespace, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cosmo_federated_graph.test", "name", name),
					resource.TestCheckResourceAttr("data.cosmo_federated_graph.test", "namespace", namespace),
				),
			},
			{
				ResourceName: "data.cosmo_federated_graph.test",
				RefreshState: true,
			},
			{
				Config:  testAccFederatedGraphDataSourceConfig(namespace, name),
				Destroy: true,
			},
		},
	})
}

func testAccFederatedGraphDataSourceConfig(namespace, name string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name      = "%s"
  namespace = cosmo_namespace.test.name
  routing_url = "https://example.com"
}

data "cosmo_federated_graph" "test" {
  name      = cosmo_federated_graph.test.name
  namespace = cosmo_federated_graph.test.namespace
}
`, namespace, name)
}
