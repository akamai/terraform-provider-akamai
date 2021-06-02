---
layout: "akamai"
page_title: "Akamai: NetworkList Network List"
subcategory: "Network Lists"
description: |-
 NetworkList Network List
---

# akamai_networklist_network_list

Use the `akamai_networklist_network_list` resource to create a network list, or to modify an existing list.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_networklist_network_list" "network_list" {
  name = "Test-white-list-NL"
  type = "IP"
  description = "network list description"
  list = var.list
  mode = "APPEND"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to be assigned to the network list.

* `type` - (Required) The type of the network list; must be either "IP" or "GEO".

* `description` - (Required) The description to be assigned to the network list.

* `list` : (Optional) A list of IP addresses or locations to be included in the list, added to an existing list, or
  removed from an existing list.

* `mode` - (Required) A string specifying the interpretation of the `list` parameter. Must be one of the following:

  * APPEND - the addresses or locations listed in `list` will be added to the network list
  * REPLACE - the addresses or locations listed in `list` will overwrite the current contents of the network list
  * REMOVE - the addresses or locations listed in `list` will be removed from the network list

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `network_list_id` - The ID of the network list.

* `sync_point` - An integer that identifies the current version of the network list; this value is incremented each time
  the list is modified. 

