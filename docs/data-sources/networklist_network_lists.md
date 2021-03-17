---
layout: "akamai"
page_title: "Akamai: NetworkLists"
subcategory: "Network Lists"
description: |-
 NetworkLists
---

# akamai_networklist_network_lists

Use the `akamai_networklist_network_lists` data source to retrieve information about the available network lists,
optionally filtered by list type or based on a search string. The information available is described
[here](https://developer.akamai.com/api/cloud_security/network_lists/v2.html#getlists). 

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_networklist_network_lists" "network_lists" {
}

//tabular data of id, name, type, elementCount, syncPoint, readonly
output "network_lists_text" {
  value = data.akamai_networklist_network_lists.network_lists.output_text
}

output "network_lists_json" {
  value = data.akamai_networklist_network_lists.network_lists.json
}

//custom output of network list ids
output "network_lists_list" {
  value = data.akamai_networklist_network_lists.network_lists.list
}

data "akamai_networklist_network_lists" "network_lists_filter" {
  name = "Test Whitelist"
  type = "IP"
}

output "network_lists_filter_text" {
  value = data.akamai_networklist_network_lists.network_lists_filter.output_text
}

output "network_lists_filter_json" {
  value = data.akamai_networklist_network_lists.network_lists_filter.json
}

//custom output of single network list id
output "network_lists_filter_list" {
  value = data.akamai_networklist_network_lists.network_lists_filter.list
}
```

## Argument Reference

* `name` - (Optional) The name of a specific network list to retrieve. If not supplied, information about all network
  lists will be returned.

* `type` - (Optional) The type of network lists to be retrieved; must be either "IP" or "GEO". If not supplied,
  information about both types will be returned.

The following arguments are supported:

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `uniqueid` - The ID of the indicated list (if the `name` argument was supplied).

* `json` - A JSON-formatted list of information about the specified network list(s).

* `output_text` - A tabular display showing the network list information.

* `list` - A list containing the IDs of the specified network lists(s).

