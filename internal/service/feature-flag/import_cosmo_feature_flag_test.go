package feature_flag_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccImportFeatureFlag(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")
	ffName := acctest.RandomWithPrefix("test-feature-flag")

	fgName, sgName, fsgName :=
		acctest.RandomWithPrefix("test-feature-flag"),
		acctest.RandomWithPrefix("test-subgraph"),
		acctest.RandomWithPrefix("test-feature-subgraph")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFeatureFlagResourceConfig(namespace, fgName, sgName, fsgName, ffName, false),
			},
			{
				ResourceName: "cosmo_feature_flag.test",
				// We can only import by namespace and name as we can't determine the ID
				ImportStateId:     namespace + "." + ffName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
