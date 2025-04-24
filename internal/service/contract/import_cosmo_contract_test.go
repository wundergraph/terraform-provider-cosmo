package contract_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccCosmoContractImportBasic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-contract")
	namespace := acctest.RandomWithPrefix("test-namespace")

	readme := "Initial readme content"

	federatedGraphName := acctest.RandomWithPrefix("test-federated-graph")

	subgraphName := acctest.RandomWithPrefix("test-subgraph")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContractResourceConfig(namespace, federatedGraphName, subgraphName, name, "internal", &readme),
			},
			{
				ResourceName: "cosmo_contract.test",

				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
