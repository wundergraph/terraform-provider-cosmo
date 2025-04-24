package namespace_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccCosmoNamespaceImportBasic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-namespace")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNamespaceResourceConfig(name),
			},
			{
				// import via ID
				ResourceName:      "cosmo_namespace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// import via name
				ResourceName:      "cosmo_namespace.test",
				ImportState:       true,
				ImportStateId:     name,
				ImportStateVerify: true,
			},
		},
	})
}
