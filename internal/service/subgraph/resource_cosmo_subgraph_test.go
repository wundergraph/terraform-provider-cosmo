package subgraph_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/api"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/wundergraph/cosmo/terraform-provider-cosmo/internal/acceptance"
)

func TestAccSubgraphResource(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")

	federatedGraphName := acctest.RandomWithPrefix("test-subgraph")
	federatedGraphRoutingURL := "https://federated-graph-example.com"

	subgraphName := acctest.RandomWithPrefix("test-subgraph")

	routingURL := "https://subgraph-example.com"
	updatedRoutingURL := "https://updated-subgraph-example.com"

	subgraphSchema := acceptance.TestAccValidSubgraphSchema
	readme := "Initial readme content"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, routingURL, subgraphSchema, readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.stage", "dev"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "subscription_protocol", api.GraphQLSubscriptionProtocolWS),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "websocket_subprotocol", api.GraphQLWebsocketSubprotocolDefault),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "readme", readme),
				),
			},
			{
				Config: testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, updatedRoutingURL, subgraphSchema, readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", updatedRoutingURL),
				),
			},
			{
				ResourceName: "cosmo_subgraph.test",
				RefreshState: true,
			},
			{
				Config:  testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, routingURL, subgraphSchema, readme),
				Destroy: true,
			},
		},
	})
}

func TestAccStandaloneSubgraphResource(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")

	federatedGraphName := acctest.RandomWithPrefix("test-subgraph")
	federatedGraphRoutingURL := "https://federated-graph-standalone-subgraph-example.com"

	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	routingURL := "https://subgraph-standalone-example.com"
	updatedRoutingURL := "https://updated-subgraph-standalone-example.com"
	subgraphSchema := acceptance.TestAccValidSubgraphSchema

	readme := "Initial readme content"
	updatedReadme := "Updated readme content"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, routingURL, subgraphSchema, readme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.stage", "dev"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "readme", readme),
				),
			},
			{
				Config: testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, updatedRoutingURL, subgraphSchema, updatedReadme),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", updatedRoutingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.stage", "dev"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "readme", updatedReadme),
				),
			},
			{
				Config: testStandaloneSubgraph(namespace, subgraphName, routingURL, subgraphSchema, nil, nil, nil, nil, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", routingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.stage", "dev"),
				),
			},
			{
				ResourceName: "cosmo_subgraph.test",
				RefreshState: true,
			},
			{
				Config:  testStandaloneSubgraph(namespace, subgraphName, routingURL, subgraphSchema, nil, nil, nil, nil, false),
				Destroy: true,
			},
		},
	})
}

func TestAccSubgraphResourceInvalidSchema(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	subgraphRoutingURL := "https://subgraph-invalid-schema-example.com"

	federatedGraphName := acctest.RandomWithPrefix("test-subgraph")
	federatedGraphRoutingURL := "https://federated-graph-invalid-subgraph-schema-example.com"
	subgraphSchema := "invalid"
	readme := "Initial readme content"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphRoutingURL, subgraphName, subgraphRoutingURL, subgraphSchema, readme),
				ExpectError: regexp.MustCompile(`.*ERR_INVALID_SUBGRAPH_SCHEMA*`),
			},
		},
	})
}

func TestAccStandaloneSubgraphResourcePublishSchema(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	subgraphRoutingURL := "https://subgraph-publish-schema-example.com"

	subgraphSchema := acceptance.TestAccValidSubgraphSchema
	updatedSubgraphSchema := "invalid"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testStandaloneSubgraph(namespace, subgraphName, subgraphRoutingURL, subgraphSchema, nil, nil, nil, nil, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", subgraphRoutingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.stage", "dev"),
				),
			},
			{
				Config:      testStandaloneSubgraph(namespace, subgraphName, subgraphRoutingURL, updatedSubgraphSchema, nil, nil, nil, nil, false),
				ExpectError: regexp.MustCompile(`.*ERR_INVALID_SUBGRAPH_SCHEMA*`),
			},
		},
	})
}

func TestOptionalValuesOfSubgraphResource(t *testing.T) {
	namespace := acctest.RandomWithPrefix("test-namespace")
	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	subgraphRoutingURL := "https://subgraph-publish-schema-example.com"

	subgraphSchema := acceptance.TestAccValidSubgraphSchema

	readme := "Initial readme content"
	subgraphSubscriptionURL := "https://subgraph-publish-schema-example.com/ws"
	subscriptionProtocol := api.GraphQLSubscriptionProtocolSSE
	websocketSubprotocol := api.GraphQLWebsocketSubprotocolGraphQLWS

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testStandaloneSubgraph(namespace, subgraphName, subgraphRoutingURL, subgraphSchema, &readme, &subgraphSubscriptionURL, &subscriptionProtocol, &websocketSubprotocol, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", subgraphRoutingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "readme", readme),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "subscription_url", subgraphSubscriptionURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "subscription_protocol", subscriptionProtocol),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "websocket_subprotocol", websocketSubprotocol),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "websocket_subprotocol", websocketSubprotocol),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.%", "0"),
				),
			},
			{
				Config: testStandaloneSubgraph(namespace, subgraphName, subgraphRoutingURL, subgraphSchema, &readme, &subgraphSubscriptionURL, &subscriptionProtocol, &websocketSubprotocol, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", subgraphRoutingURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "readme", readme),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "subscription_url", subgraphSubscriptionURL),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "subscription_protocol", subscriptionProtocol),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "websocket_subprotocol", websocketSubprotocol),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.team", "backend"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.stage", "dev"),
				),
			},
			{
				Config: testStandaloneSubgraph(namespace, subgraphName, subgraphRoutingURL, subgraphSchema, nil, nil, nil, nil, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "name", subgraphName),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", namespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "routing_url", subgraphRoutingURL),
					resource.TestCheckNoResourceAttr("cosmo_subgraph.test", "readme"),
					resource.TestCheckNoResourceAttr("cosmo_subgraph.test", "subscription_url"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "labels.%", "0"),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "subscription_protocol", api.GraphQLSubscriptionProtocolWS),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "websocket_subprotocol", api.GraphQLWebsocketSubprotocolDefault),
				),
			},
		},
	})
}

func TestAccSubgraphResourceNamespaceChangeForceNew(t *testing.T) {
	initialNamespace := acctest.RandomWithPrefix("test-namespace-1")
	newNamespace := acctest.RandomWithPrefix("test-namespace-2")

	subgraphName := acctest.RandomWithPrefix("test-subgraph")
	routingURL := "https://subgraph-example.com"
	subgraphSchema := acceptance.TestAccValidSubgraphSchema

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Initial creation of two namespaces
			{
				Config: testStandaloneSubgraph(initialNamespace, subgraphName, routingURL, subgraphSchema, nil, nil, nil, nil, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_namespace.test", "name", initialNamespace),
				),
			},
			{
				Config: testSubgraphWithTwoNamespaces(initialNamespace, newNamespace, subgraphName, routingURL, subgraphSchema, nil, nil, nil, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cosmo_namespace.test2", "name", newNamespace),
					resource.TestCheckResourceAttr("cosmo_subgraph.test", "namespace", newNamespace),
				),
			},
		},
	})
}

func testAccSubgraphResourceConfig(namespace, federatedGraphName, federatedGraphroutingURL, subgraphName, subgraphRoutingURL, subgraphSchema, readme string) string {
	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_federated_graph" "test" {
  name      	= "%s"
  namespace 	= cosmo_namespace.test.name
  routing_url 	= "%s"
  label_matchers = ["team=backend"]

  depends_on = [cosmo_subgraph.test]
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "%s"
  schema              = <<-EOT
  %s
  EOT
  labels              = { 
  	"team"	= "backend", 
	"stage" = "dev" 
  }
  readme              =  "%s"
}
`, namespace, federatedGraphName, federatedGraphroutingURL, subgraphName, subgraphRoutingURL, subgraphSchema, readme)
}

func testStandaloneSubgraph(namespace, subgraphName, subgraphRoutingURL, subgraphSchema string, readme, subscriptionUrl, subscriptionProtocol, websocketSubprotocol *string, unsetLabels bool) string {
	var readmePart, subscriptionUrlPart, subscriptionProtocolPart, websocketSubprotocolPart string
	if readme != nil {
		readmePart = fmt.Sprintf(`readme = "%s"`, *readme)
	}

	if subscriptionUrl != nil {
		subscriptionUrlPart = fmt.Sprintf(`subscription_url = "%s"`, *subscriptionUrl)
	}

	if subscriptionProtocol != nil {
		subscriptionProtocolPart = fmt.Sprintf(`subscription_protocol = "%s"`, *subscriptionProtocol)
	}

	if websocketSubprotocol != nil {
		websocketSubprotocolPart = fmt.Sprintf(`websocket_subprotocol = "%s"`, *websocketSubprotocol)
	}

	if unsetLabels {
		return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "%s"
  schema              = <<-EOT
  %s
  EOT
  %s
  %s
  %s
  %s
}
`, namespace, subgraphName, subgraphRoutingURL, subgraphSchema, readmePart, subscriptionUrlPart, subscriptionProtocolPart, websocketSubprotocolPart)
	}

	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test.name
  routing_url         = "%s"
  schema              = <<-EOT
  %s
  EOT
  labels              = { 
  	"team"	= "backend", 
	"stage" = "dev" 
  }
  %s
  %s
  %s
  %s
}
`, namespace, subgraphName, subgraphRoutingURL, subgraphSchema, readmePart, subscriptionUrlPart, subscriptionProtocolPart, websocketSubprotocolPart)
}

func testSubgraphWithTwoNamespaces(oldNamespace, newNamespace, subgraphName, subgraphRoutingURL, subgraphSchema string, readme, subscriptionUrl, subscriptionProtocol, websocketSubprotocol *string) string {
	var readmePart, subscriptionUrlPart, subscriptionProtocolPart, websocketSubprotocolPart string
	if readme != nil {
		readmePart = fmt.Sprintf(`readme = "%s"`, *readme)
	}

	if subscriptionUrl != nil {
		subscriptionUrlPart = fmt.Sprintf(`subscription_url = "%s"`, *subscriptionUrl)
	}

	if subscriptionProtocol != nil {
		subscriptionProtocolPart = fmt.Sprintf(`subscription_protocol = "%s"`, *subscriptionProtocol)
	}

	if websocketSubprotocol != nil {
		websocketSubprotocolPart = fmt.Sprintf(`websocket_subprotocol = "%s"`, *websocketSubprotocol)
	}

	return fmt.Sprintf(`
resource "cosmo_namespace" "test" {
  name = "%s"
}

resource "cosmo_namespace" "test2" {
  name = "%s"
}

resource "cosmo_subgraph" "test" {
  name                = "%s"
  namespace           = cosmo_namespace.test2.name
  routing_url         = "%s"
  schema              = <<-EOT
  %s
  EOT
  %s
  %s
  %s
  %s
}
`, oldNamespace, newNamespace, subgraphName, subgraphRoutingURL, subgraphSchema, readmePart, subscriptionUrlPart, subscriptionProtocolPart, websocketSubprotocolPart)
}
