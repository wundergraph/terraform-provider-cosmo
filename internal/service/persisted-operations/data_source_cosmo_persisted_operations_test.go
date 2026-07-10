package persisted_operations_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccPersistedOperationsDataSource(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")
	graphName := acctest.RandomWithPrefix("test-graph")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	clientName := acctest.RandomWithPrefix("test-client")

	oneOperation := fmt.Sprintf(`{ capsules = %q }`, testAccCapsulesOperation)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPersistedOperationsDataSourceConfig(namespace, graphName, subgraphName, clientName, oneOperation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.cosmo_persisted_operations.test", "client_name", clientName),
					resource.TestCheckResourceAttr("data.cosmo_persisted_operations.test", "operations.#", "1"),
					resource.TestCheckResourceAttrSet("data.cosmo_persisted_operations.test", "client_id"),
					resource.TestCheckResourceAttrSet("data.cosmo_persisted_operations.test", "operations.0.id"),
					resource.TestCheckResourceAttr("data.cosmo_persisted_operations.test", "operations.0.contents", testAccCapsulesOperation),
				),
			},
		},
	})
}

func testAccPersistedOperationsDataSourceConfig(namespace, graphName, subgraphName, clientName, operations string) string {
	return fmt.Sprintf(`
%s

data "cosmo_persisted_operations" "test" {
  federated_graph_name = cosmo_federated_graph.test.name
  namespace            = cosmo_namespace.test.name
  client_name          = cosmo_persisted_operations.test.client_name
}
`, testAccPersistedOperationsResourceConfig(namespace, graphName, subgraphName, clientName, operations))
}
