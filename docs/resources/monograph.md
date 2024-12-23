---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cosmo_monograph Resource - cosmo"
subcategory: ""
description: |-
  A monograph is a resource that represents a single subgraph with GraphQL Federation disabled.
  For more information on monographs, please refer to the Cosmo Documentation https://cosmo-docs.wundergraph.com/cli/monograph.
---

# cosmo_monograph (Resource)

A monograph is a resource that represents a single subgraph with GraphQL Federation disabled.

For more information on monographs, please refer to the [Cosmo Documentation](https://cosmo-docs.wundergraph.com/cli/monograph).

## Example Usage

```terraform
resource "cosmo_monograph" "example" {
  name        = var.monograph_name
  namespace   = var.monograph_namespace
  graph_url   = var.monograph_graph_url
  routing_url = var.monograph_routing_url
  schema      = var.monograph_schema
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `graph_url` (String) The GraphQL endpoint URL of the monograph.
- `name` (String) The name of the monograph.
- `routing_url` (String) The routing URL for the monograph.

### Optional

- `admission_webhook_secret` (String) The admission webhook secret for the monograph.
- `admission_webhook_url` (String) The admission webhook URL for the monograph.
- `namespace` (String) The namespace in which the monograph is located.
- `readme` (String) The readme for the subgraph.
- `schema` (String) The schema for the subgraph.
- `subscription_protocol` (String) The subscription protocol for the subgraph.
- `subscription_url` (String) The subscription URL for the subgraph.
- `websocket_subprotocol` (String) The websocket subprotocol for the subgraph.

### Read-Only

- `id` (String) The unique identifier of the monograph resource.
