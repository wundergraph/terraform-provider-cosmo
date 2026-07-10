package persisted_operations_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"
)

const (
	testAccCapsulesOperation = "query Capsules { capsules { id } }"
	testAccDragonsOperation  = "query Dragons { dragons { id } }"
)

func TestAccPersistedOperationsResource(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")
	graphName := acctest.RandomWithPrefix("test-graph")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	clientName := acctest.RandomWithPrefix("test-client")

	oneOperation := fmt.Sprintf(`{ capsules = %q }`, testAccCapsulesOperation)
	twoOperations := fmt.Sprintf(`{ capsules = %q, dragons = %q }`, testAccCapsulesOperation, testAccDragonsOperation)
	otherOperation := fmt.Sprintf(`{ dragons = %q }`, testAccDragonsOperation)

	checkRemoteOperationCount := func(expected int) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			apiClient, err := api.NewClient(os.Getenv("COSMO_API_KEY"), os.Getenv("COSMO_API_URL"))
			if err != nil {
				return errors.New("Error creating api client")
			}
			client, apiError := apiClient.GetClient(context.Background(), graphName, namespace, clientName)
			if apiError != nil {
				return fmt.Errorf("Error getting client: %s", apiError.Error())
			}
			operations, apiError := apiClient.GetPersistedOperations(context.Background(), graphName, namespace, client.Id)
			if apiError != nil {
				return fmt.Errorf("Error getting persisted operations: %s", apiError.Error())
			}
			if len(operations) != expected {
				return fmt.Errorf("Expected %d persisted operations, got %d", expected, len(operations))
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPersistedOperationsResourceConfig(namespace, graphName, subgraphName, clientName, oneOperation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_persisted_operations.test", "federated_graph_name", graphName),
					resource.TestCheckResourceAttr("cosmo_persisted_operations.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_persisted_operations.test", "client_name", clientName),
					resource.TestCheckResourceAttr("cosmo_persisted_operations.test", "operations.%", "1"),
					resource.TestCheckResourceAttrSet("cosmo_persisted_operations.test", "id"),
				),
			},
			{
				Config: testAccPersistedOperationsResourceConfig(namespace, graphName, subgraphName, clientName, oneOperation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_persisted_operations.test", "operations.%", "1"),
				),
			},
			{
				ResourceName: "cosmo_persisted_operations.test",
				RefreshState: true,
			},
			{
				Config: testAccPersistedOperationsResourceConfig(namespace, graphName, subgraphName, clientName, twoOperations),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_persisted_operations.test", "operations.%", "2"),
					checkRemoteOperationCount(2),
				),
			},
			{
				Config: testAccPersistedOperationsResourceConfig(namespace, graphName, subgraphName, clientName, otherOperation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_persisted_operations.test", "operations.%", "1"),
					checkRemoteOperationCount(1),
				),
			},
			{
				Config: testAccNoPersistedOperationsResourceConfig(namespace, graphName, subgraphName),
				Check: resource.ComposeTestCheckFunc(
					checkRemoteOperationCount(0),
				),
			},
		},
	})
}

func TestAccPersistedOperationsResourceClientNameRecreates(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")
	graphName := acctest.RandomWithPrefix("test-graph")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	clientName := acctest.RandomWithPrefix("test-client")
	newClientName := acctest.RandomWithPrefix("test-client-new")

	oneOperation := fmt.Sprintf(`{ capsules = %q }`, testAccCapsulesOperation)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPersistedOperationsResourceConfig(namespace, graphName, subgraphName, clientName, oneOperation),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_persisted_operations.test", "client_name", clientName),
				),
			},
			{
				Config: testAccPersistedOperationsResourceConfig(namespace, graphName, subgraphName, newClientName, oneOperation),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("cosmo_persisted_operations.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_persisted_operations.test", "client_name", newClientName),
				),
			},
		},
	})
}

func testAccNoPersistedOperationsResourceConfig(namespace, graphName, subgraphName string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name           = "%s"
  namespace      = cosmo_namespace.test.name
  routing_url    = "https://example.com"
  label_matchers = ["team=backend"]

  depends_on = [cosmo_subgraph.test]
}

resource "cosmo_subgraph" "test" {
  name        = "%s"
  namespace   = cosmo_namespace.test.name
  routing_url = "https://subgraph-example.com"
  schema      = <<-EOT
  %s
  EOT
  labels = {
    "team" = "backend"
  }
}
`, namespace, graphName, subgraphName, acceptance.TestAccValidSubgraphSchema)
}

func testAccPersistedOperationsResourceConfig(namespace, graphName, subgraphName, clientName, operations string) string {
	return fmt.Sprintf(`
%s

resource "cosmo_persisted_operations" "test" {
  federated_graph_name = cosmo_federated_graph.test.name
  namespace            = cosmo_namespace.test.name
  client_name          = "%s"
  operations           = %s
}
`, testAccNoPersistedOperationsResourceConfig(namespace, graphName, subgraphName), clientName, operations)
}
