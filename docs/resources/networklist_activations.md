---
layout: akamai
subcategory: Network Lists
---

# akamai_networklist_activations

Use the `akamai_networklist_activations` resource to activate a network list in either the STAGING or PRODUCTION
environment.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_networklist_network_list" "network_list_ip" {
  name        = var.ip_list_name
  type        = "IP"
  description = "IP network list"
  list        = var.ip_list
  mode        = "REPLACE"
}

resource "akamai_networklist_activations" "activation" {
  network_list_id     = resource.akamai_networklist_network_list.network_list_ip.network_list_id
  network             = "STAGING"
  sync_point          = resource.akamai_networklist_network_list.network_list_ip.sync_point
  notes               = "TEST Notes"
  notification_emails = ["user@example.com"]
}
```

## Argument Reference

The following arguments are supported:

* `network_list_id` - (Required) The ID of the network list to be activated

* `network` - (Optional) The network to be used, either `STAGING` or `PRODUCTION`. If not supplied, defaults to
  `STAGING`.

* `sync_point` - (Required) An integer that identifies the current version of the network list; this value is incremented each time
  the list is modified.

* `notes` - (Optional) A comment describing the activation.

* `notification_emails` - (Required) A bracketed, comma-separated list of email addresses that will be notified when the
  operation is complete.

## Attributes Reference

In addition to the arguments above, the following attribute is exported:

* `status` - The string `ACTIVATED` if the activation was successful, or a string identifying the reason why the network
  list was not activated.

