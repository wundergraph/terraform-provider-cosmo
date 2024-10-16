# Cosmo Terraform Provider

This repository is for the [Cosmo][https://registry.terraform.io/wundergraph/cosmo](https://registry.terraform.io/providers/wundergraph/cosmo/latest/docs) Terraform provider, designed to manage Cosmo resources within Terraform. It includes a resource and a data source, examples, and generated documentation.

## Requirements

- [Terraform](https://developer.hashiCorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Implemented Resources

The Cosmo Terraform provider includes the following resources and data sources:

### Resources

- [cosmo_namespace](docs/resources/namespace.md): Manages namespaces within Cosmo.
- [cosmo_monograph](docs/resources/monograph.md): Manages monographs in Cosmo.
- [cosmo_federated_graph](docs/resources/federated_graph.md): Manages federated graphs in Cosmo.
- [cosmo_subgraph](docs/resources/subgraph.md): Manages subgraphs in Cosmo.
- [cosmo_router_token](docs/resources/cosmo_router_token.md): Retrieves information about subgraphs in Cosmo.

### Data Sources

- [cosmo_namespace](docs/data-sources/namespace.md): Retrieves information about namespaces in Cosmo.
- [cosmo_monograph](docs/data-sources/monograph.md): Retrieves information about monographs in Cosmo.
- [cosmo_federated_graph](docs/data-sources/federated_graph.md): Retrieves information about federated graphs in Cosmo.
- [cosmo_subgraph](docs/data-sources/subgraph.md): Retrieves information about subgraphs in Cosmo.

Each resource and data source allows you to define and manage specific aspects of your Cosmo infrastructure seamlessly within Terraform.

## Example Usage

The provider can be used as follows:

```hcl
# variables.tf
variable "namespace" {
  type        = string
  description = "The name of the namespace to be used for the federated graph"
  default = "your-namespace"
}

variable "federated_graph" {
  type = object({
    name           = string
    routing_url    = string
    label_matchers = list(string)
  })
  description = "The parameters of the federated graph"
  default = {
    name = "your-federated-graph"
    routing_url = "http://localhost:3000"
    label_matchers = ["team=backend", "stage=dev"]
  }
}

variable "subgraphs" {
  type = map(object({
    name        = string
    routing_url = string
    labels      = map(string)
    schema      = string
  }))
  description = "The subgraphs to be added to the federated graph"
  default = {
    "subgraph-1" = {
      name = "your-subgraph-1"
      routing_url = "http://example.com/routing"
      schema = "type Query { hello: String }"
      labels = {
        "team" = "backend"
        "stage" = "dev"
      }
    }
  }
}

variable "router_token_name" {
  type        = string
  description = "The name of the router token to be created"
  default = "your-router-token"
}

# main.tf
resource "cosmo_namespace" "namespace" {
  name = var.namespace
}

resource "cosmo_federated_graph" "federated_graph" {
  name           = var.federated_graph.name
  routing_url    = var.federated_graph.routing_url
  namespace      = cosmo_namespace.namespace.name
  label_matchers = var.federated_graph.label_matchers

  depends_on     = [cosmo_subgraph.subgraph]
}

resource "cosmo_subgraph" "subgraph" {
  for_each = var.subgraphs

  name        = each.value.name
  namespace   = cosmo_namespace.namespace.name
  routing_url = each.value.routing_url
}

resource "cosmo_router_token" "router_token" {
  name       = var.router_token_name
  namespace  = cosmo_namespace.namespace.name
  graph_name = cosmo_federated_graph.federated_graph.name
}

# outputs.tf
output "router_token" {
  value = cosmo_router_token.router_token.token
}
```

Further in depth examples can be found in the [examples](examples) directory.

## Cosmo Local Example

The module [cosmo-local](examples/cosmo-local) contains an example of how to use the provider to manage a local cosmo setup on minikube.

It will create a minikube cluster, install cosmo and other dependencies and also setup a federated graph with a subgraph and deploy a router with a router token.

To run the example, run `make e2e-apply-cosmo-local` from the root of the repository.

Running apply will print out the hosts you need to add to your local `/etc/hosts` file to access the services:

```
# example output
hosts = <<EOT
    # WunderGraph
    192.168.49.2 studio.wundergraph.local
    192.168.49.2 controlplane.wundergraph.local
    192.168.49.2 router.wundergraph.local
    192.168.49.2 keycloak.wundergraph.local
    192.168.49.2 otelcollector.wundergraph.local
    192.168.49.2 graphqlmetrics.wundergraph.local
    192.168.49.2 cdn.wundergraph.local
EOT
```

You can now access the router on `router.wundergraph.local` and the studio on `studio.wundergraph.local`. To test your installation.

## Building The Provider

To build the provider, clone the repository, enter the directory, and run `make install` to compile and install the provider binary. Note that the `install` command will first build the provider to ensure the binary is up to date.

## Usage

Ensure to set `COSMO_API_URL` and `COSMO_API_KEY` environment variables to point to your cosmo setup.

For example:

```bash
export COSMO_API_KEY="<cosmo-api-token>"
export COSMO_API_URL="http://localhost:3001"


# start cosmo from within the cosmo repo
cd cosmo
make full-demo-up

# build install and run the e2e tests with the cosmo provider
cd terraform-provider-cosmo
make clean build install e2e
```

The following commands are used to build and install the provider binary locally for use with end-to-end tests:

1. **Install the Provider**: Run the following command to build and install the provider binary locally for use with end-to-end tests:

   ```bash
   make clean build install
   ```

2. **Run Tests**: Execute acceptance tests to ensure the provider works as expected:

   ```bash
   make testacc
   ```

3. **Generate Files**: Update any generated files with this command:

   ```bash
   make generate
   ```

4. **Format Code**: Format Go and Terraform files for consistency:

   ```bash
   make fmt
   ```

5. **Build for All Architectures**: Compile the provider for various operating systems and architectures:

   ```bash
   make build-all-arches
   ```

## Makefile Tasks

The Makefile includes several tasks to facilitate development and testing. For local development, `make build install` should be used to install the provider locally.

### General Build Tasks

- **default**: Runs acceptance tests.
- **testacc**: Runs tests with a timeout.
- **test-go**: Runs Go tests.
- **test**: Cleans, builds, installs, runs acceptance tests, and executes end-to-end tests.
- **generate**: Updates generated files.
- **tidy**: Cleans up the `go.mod` file.
- **fmt**: Formats Go and Terraform files for consistency.
- **build**: Compiles the provider binary.
- **install**: Installs the binary in the Terraform plugin directory after building it.
- **clean-local**: Cleans up local build artifacts.
- **build-all-arches**: Compiles the binary for multiple OS and architectures.

### End-to-End (E2E) Tasks

- **e2e-apply-cd**: Runs end-to-end tests for the CD feature (points to `examples/provider`).
- **e2e-destroy-cd**: Cleans up after CD tests (points to `examples/provider`).
- **e2e-clean-cd**: Cleans up CD test artifacts (points to `examples/provider`).
- **e2e-apply-cosmo**: Runs end-to-end tests for the Cosmo feature (points to `examples/guides/cosmo`).
- **e2e-destroy-cosmo**: Cleans up after Cosmo tests (points to `examples/guides/cosmo`).
- **e2e-clean-cosmo**: Cleans up Cosmo test artifacts (points to `examples/guides/cosmo`).
- **e2e-apply-cosmo-monograph**: Runs end-to-end tests for the Cosmo monograph feature (points to `examples/guides/cosmo-monograph`).
- **e2e-destroy-cosmo-monograph**: Cleans up after Cosmo monograph tests (points to `examples/guides/cosmo-monograph`).
- **e2e-clean-cosmo-monograph**: Cleans up Cosmo monograph test artifacts (points to `examples/guides/cosmo-monograph`).
- **e2e-apply-cosmo-monograph-contract**: Runs end-to-end tests for the Cosmo monograph contract feature (points to `examples/guides/cosmo-monograph-contract`).
- **e2e-destroy-cosmo-monograph-contract**: Cleans up after Cosmo monograph contract tests (points to `examples/guides/cosmo-monograph-contract`).
- **e2e-clean-cosmo-monograph-contract**: Cleans up Cosmo monograph contract test artifacts (points to `examples/guides/cosmo-monograph-contract`).
- **e2e-apply-cosmo-local**: Runs end-to-end tests for the local Cosmo setup (points to `examples/guides/cosmo-local`).
- **e2e-destroy-cosmo-local**: Cleans up after local Cosmo tests (points to `examples/guides/cosmo-local`).
- **e2e-clean-cosmo-local**: Cleans up local Cosmo test artifacts (points to `examples/guides/cosmo-local`).
- **e2e-cd**: Runs both apply and destroy for CD tests.
- **e2e-cosmo**: Runs both apply and destroy for Cosmo tests.
- **e2e-cosmo-monograph**: Runs both apply and destroy for Cosmo monograph tests.
- **e2e-cosmo-monograph-contract**: Runs both apply and destroy for Cosmo monograph contract tests.
- **e2e-cosmo-local**: Runs both apply and destroy for local Cosmo tests.
- **e2e**: Runs all end-to-end tests.
- **clean**: Cleans up all test artifacts and local builds.
- **destroy**: Cleans up all resources created by the tests.

## Releasing

The Terraform Provider can be release by triggering the `Release` workflow in the `.github/workflows` directory.
This workflow must be triggered manually on the main branch when a release is needed.

The workflow will create a new tag and push it to the remote.
Subsequently the workflow will build the go release and create a new github release.

Tags follow this schema: `vX.Y.Z`. This is needed for the terraform registry to pick up new versions.
