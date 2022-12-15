---
layout: akamai
subcategory: Property Provisioning
---

# akamai_property_include (Beta)

Use the `akamai_property_include` data source to get details about a specific include.

## Basic usage

This example returns details for an include based on contract, group, and include IDs.

```hcl
data "akamai_property_include" "my_example" {
  contract_id = "ctr_1234"
  group_id    = "grp_5678"
  include_id  = "inc_9012"
}

output "my_example" {
  value = data.akamai_property_include.my_example
}
```

## Argument reference

This data source supports these arguments:

* `contract_id` - (Required) A contract's unique ID, including the optional `ctr_` prefix.
* `group_id` - (Required) A group's unique ID, including the optional `grp_` prefix.
* `include_id` - (Required) An include's unique ID with the optional `inc_` prefix.

## Attributes reference

This data source returns these attributes:

* `name` - The descriptive name for the include.
* `type` - Specifies the type of the include, either `MICROSERVICES` or `COMMON_SETTINGS`. Use this field for filtering. `MICROSERVICES` allow different teams to work independently on different parts of a single site. `COMMON_SETTINGS` includes are useful for configurations that share a large number of settings, often managed by a central team.
* `latest_version` - Returns the most recent version of the include.
* `staging_version` - The include version currently activated on the staging network.
* `production_version` - The include version currently activated on the production network.
