---
layout: akamai
subcategory: Network Lists
---

# akamai_networklist_description

Use the `akamai_networklist_description` resource to update the name or description of an existing network list.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_networklist_description" "network_list_description" {
  network_list_id = var.network_list_id
  name = "Test network list updated name"
  description = "Test network list updated description"
}
```

## Argument Reference

The following arguments are supported:

* `network_list_id` - (Required) The unique ID of the network list to use.

* `name` - (Required) The name to be assigned to the network list.

* `description` - (Required) The description to be assigned to the network list.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

