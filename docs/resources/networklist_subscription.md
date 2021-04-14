---
layout: "akamai"
page_title: "Akamai: NetworkLists Subscription"
subcategory: "Network Lists"
description: |-
 NetworkList Subscription
---

# akamai_networklist_subscription

Use the `akamai_networklist_subscription` resource to specify a set of email addresses to be notified of changes to any
of a set of network lists.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_networklist_network_lists" "network_lists_filter" {
  name = var.network_list
}

resource "akamai_networklist_subscription" "subscribe" {
  network_list = data.akamai_networklist_network_lists.network_lists_filter.list
  recipients = ["user@example.com"]
}
```

## Argument Reference

The following arguments are supported:

* `network_list` - (Required) A list containing one or more IDs of the network lists to which the indicated email
  addresses should be subscribed.

* `recipients` - (Required) A bracketed, comma-separated list of email addresses that will be notified of changes to any
  of the specified network lists.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

