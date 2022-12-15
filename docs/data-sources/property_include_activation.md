---
layout: akamai
subcategory: Property Provisioning
---

# akamai_property_include_activation (Beta)

Use the `akamai_property_include_activation` data source to get activation details for an include on the provided network.

## Basic usage

This example returns the include activation on a specified network based on the contract, group, and include IDs.

```hcl
data "akamai_include_activation" "my_example" {
  contract_id = "ctr_1234"
  group_id    = "grp_5678"
  include_id  = "inc_9012"
  network     = "PRODUCTION"
}

output "my_example" {
  value = data.akamai_include_activation.my_example
}
```

## Argument reference

This data source supports these arguments:

* `contract_id` - (Required) A contract's unique ID, including the optional `ctr_` prefix.
* `group_id` - (Required) A group's unique ID, including the optional `grp_` prefix.
* `include_id` - (Required) An include's unique ID with the optional `inc_` prefix.
* `network` - (Required) The Akamai network where you want to check the activation details, either `STAGING` or `PRODUCTION`. `STAGING` is the default.

## Attributes reference

This data source returns these attributes:

* `version` - The version of the activated include.
* `name` - The descriptive name for the property.
* `note` - A log message assigned to the activation request.
* `notify_emails` - The list of email addresses notified when the activation status changes.
