package persisted_operations_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccPersistedOperationsImport(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")
	graphName := acctest.RandomWithPrefix("test-graph")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	clientName := acctest.RandomWithPrefix("test-client")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPersistedOperationsResourceConfig(namespace, graphName, subgraphName, clientName, testAccOneOperation),
			},
			{
				ResourceName:      "cosmo_persisted_operations.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s:%s:%s", graphName, namespace, clientName),
				ImportStateVerify: false,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 state, got %d", len(states))
					}
					state := states[0]
					if state.Attributes["operations.%"] != "1" {
						return fmt.Errorf("expected 1 operation, got %s", state.Attributes["operations.%"])
					}
					if state.Attributes["client_name"] != clientName {
						return fmt.Errorf("expected client name %s, got %s", clientName, state.Attributes["client_name"])
					}
					if state.Attributes["id"] == "" {
						return errors.New("expected id to be set")
					}
					return nil
				},
			},
		},
	})
}
