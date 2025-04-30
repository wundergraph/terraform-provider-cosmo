package router_token_test

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

func TestAccTokenResource(t *testing.T) {
	name := acctest.RandomWithPrefix("test-token")
	namespace := acctest.RandomWithPrefix("test-namespace")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTokenResourceConfig(namespace, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_router_token.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_router_token.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_router_token.test", "graph_name", "federated-graph"),
				),
			},
			{
				Config: testAccTokenResourceConfig(namespace, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_router_token.test", "graph_name", "federated-graph"),
				),
			},
			{
				ResourceName: "cosmo_router_token.test",
				RefreshState: true,
			},
			{
				Config: testAccNoTokenResourceConfig(namespace),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						apiClient, err := api.NewClient(os.Getenv("COSMO_API_KEY"), os.Getenv("COSMO_API_URL"))
						if err != nil {
							return errors.New("Error creating api client")
						}
						_, errGetToken := apiClient.GetToken(context.Background(), name, "federated-graph", namespace)
						if errGetToken == nil {
							return errors.New("Token should not exists")
						}
						if errGetToken.Err.Error() != "ErrNotFound" {
							return fmt.Errorf("Error should be not found: %s", errGetToken.Err.Error())
						}

						return nil
					},
				),
			},
		},
	})
}

func TestAccTokenResourceUpdateRecreates(t *testing.T) {
	name := acctest.RandomWithPrefix("test-token")
	newName := acctest.RandomWithPrefix("test-token-new")
	namespace := acctest.RandomWithPrefix("test-namespace")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTokenResourceConfig(namespace, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_router_token.test", "name", name),
					resource.TestCheckResourceAttr("cosmo_router_token.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_router_token.test", "graph_name", "federated-graph"),
				),
			},
			{
				Config: testAccTokenResourceConfig(namespace, newName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("cosmo_router_token.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_router_token.test", "name", newName),
					resource.TestCheckResourceAttr("cosmo_router_token.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_router_token.test", "graph_name", "federated-graph"),
					func(s *terraform.State) error {
						apiClient, err := api.NewClient(os.Getenv("COSMO_API_KEY"), os.Getenv("COSMO_API_URL"))
						if err != nil {
							return errors.New("Error creating api client")
						}
						_, errGetToken := apiClient.GetToken(context.Background(), name, "federated-graph", namespace)
						if errGetToken == nil {
							return errors.New("Token should not exists")
						}
						if errGetToken.Err.Error() != "ErrNotFound" {
							return fmt.Errorf("Error should be not found: %s", errGetToken.Err.Error())
						}

						return nil
					},
				),
			},
		},
	})
}

func testAccNoTokenResourceConfig(namespace string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name      	= "federated-graph"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "https://example.com"
  readme    	= "This is a test federated graph"
}
`, namespace)
}

func testAccTokenResourceConfig(namespace, name string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name      	= "federated-graph"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "https://example.com"
  readme    	= "This is a test federated graph"
}

resource "cosmo_router_token" "test" {
  name       = "%s"
  namespace  = cosmo_namespace.test.name
  graph_name = cosmo_federated_graph.test.name
}
`, namespace, name)
}
