---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cosmo_namespace Data Source - cosmo"
subcategory: ""
description: |-
  Cosmo Namespace Data Source
---

# cosmo_namespace (Data Source)

Cosmo Namespace Data Source

## Example Usage

```terraform
data "cosmo_namespace" "test" {
  name = var.name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the namespace.

### Read-Only

- `id` (String) The unique identifier of the namespace resource, automatically generated by the system.
