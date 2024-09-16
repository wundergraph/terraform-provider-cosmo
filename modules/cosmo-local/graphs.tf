resource "random_string" "module_prefix" {
  length  = 6
  special = false
}

locals {
  prefix = lower(random_string.module_prefix.result)
}

module "cosmo_federated_graph" {
  source            = "../cosmo-federated-graph"
  namespace         = "${var.stage}-${local.prefix}"
  router_token_name = "${var.stage}-${local.prefix}-router-token"

  federated_graph = {
    name        = "${local.prefix}-federated-graph"
    routing_url = "http://localhost:3000"
    label_matchers = [
      "team=backend",
    ]
  }
  subgraphs = {
    "availability" = {
      name        = "${var.stage}-${local.prefix}-availability"
      routing_url = "http://example.com/routing"
      schema      = <<EOF
      extend schema
      @link(url: "https://specs.apollo.dev/federation/v2.5", import: ["@authenticated", "@composeDirective", "@external", "@extends", "@inaccessible", "@interfaceObject", "@override", "@provides", "@key", "@requires", "@requiresScopes", "@shareable", "@tag"])

      type Mutation {
          updateAvailability(employeeID: Int!, isAvailable: Boolean!): Employee!
      }
      type Employee @key(fields: "id") {
        id: Int!
        isAvailable: Boolean!
      }
      EOF
      readme      = "Availabilty Subgraph"
      labels = {
        "team"  = "backend"
        "stage" = "${var.stage}"
      }
    }
    "products" = {
      name        = "${var.stage}-${local.prefix}-products"
      routing_url = "http://example.com/routing"
      schema      = <<EOF
      extend schema
      @link(url: "https://specs.apollo.dev/federation/v2.5", import: ["@authenticated", "@composeDirective", "@external", "@extends", "@inaccessible", "@interfaceObject", "@override", "@provides", "@key", "@requires", "@requiresScopes", "@shareable", "@tag"])

      schema {
        query: Queries
        mutation: Mutation
      }

      type Queries {
        productTypes: [Products!]!
        topSecretFederationFacts: [TopSecretFact!]! @requiresScopes(scopes: [["read:fact"], ["read:all"]])
        factTypes: [TopSecretFactType!]
      }

      type Mutation {
        addFact(fact: TopSecretFactInput!): TopSecretFact! @requiresScopes(scopes: [["write:fact"], ["write:all"]])
      }

      input TopSecretFactInput {
        title: String!
        description: FactContent!
        factType: TopSecretFactType!
      }

      enum TopSecretFactType @authenticated {
        DIRECTIVE,
        ENTITY,
        MISCELLANEOUS,
      }

      interface TopSecretFact @authenticated {
        description: FactContent!
        factType: TopSecretFactType
      }

      scalar FactContent @requiresScopes(scopes: [["read:scalar"], ["read:all"]])

      type DirectiveFact implements TopSecretFact @authenticated {
        title: String!
        description: FactContent!
        factType: TopSecretFactType
      }

      type EntityFact implements TopSecretFact @requiresScopes(scopes: [["read:entity"]]){
        title: String!
        description: FactContent!
        factType: TopSecretFactType
      }

      type MiscellaneousFact implements TopSecretFact {
        title: String!
        description: FactContent! @requiresScopes(scopes: [["read:miscellaneous"]])
        factType: TopSecretFactType
      }

      enum ProductName {
        CONSULTANCY
        COSMO
        ENGINE
        FINANCE
        HUMAN_RESOURCES
        MARKETING
        SDK
      }

      type Employee @key(fields: "id") {
        id: Int!
        products: [ProductName!]!
        notes: String @override(from: "employees")
      }

      union Products = Consultancy | Cosmo | Documentation

      type Consultancy @key(fields: "upc") {
        upc: ID!
        name: ProductName!
      }

      type Cosmo @key(fields: "upc") {
        upc: ID!
        name: ProductName!
        repositoryURL: String!
      }

      type Documentation {
        url(product: ProductName!): String!
        urls(products: [ProductName!]!): [String!]!
      }
      EOF
      readme      = "Products Subgraph"
      labels = {
        "team"  = "backend"
        "stage" = "${var.stage}"
      }
    }
  }

  depends_on = [
    module.minikube,
    module.cosmo_release
  ]
}
 