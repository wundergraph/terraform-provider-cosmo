package namespace_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
	"testing"
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
