---
layout: akamai
subcategory: Property Provisioning
---

# akamai_property_include_parents (Beta)

Use the `akamai_property_include_parents` data source to return a list of parent properties that use the given include. In your property's rule tree, you can reference an include by adding the `include` behavior and specifying the `include_id`.

## Basic usage

This example returns all active properties a specific include is referenced in, based on the contract, group, and include IDs.

```hcl
data "akamai_property_include_parents" "my_example" {
    contract_id = "ctr_1-AB123"
    group_id    = "grp_12345"
    include_id  = "inc_123456"
}

output "my_example" {
  value = data.akamai_property_include_parents.my_example
}
```

## Argument reference

This data source supports these arguments:

* `contract_id` - (Required) A contract's unique ID, including the optional `ctr_` prefix.
* `group_id` - (Required) A group's unique ID, including the optional `grp_` prefix.
* `include_id` - (Required) An include's unique ID with the optional `inc_` prefix.

## Attributes reference

This data source returns these attributes:

* `parents` - The list of include's parent properties.
 * `id` - The property's unique identifier.
 * `name` - The descriptive name for the property.
 * `staging_version` - The property version currently activated on the staging network.
 * `production_version` - The property version currently activated on the production network.
 * `is_include_used_in_staging_version` - Whether the specified include is active on the staging network and is referenced in parent's `staging_version`.
 * `is_include_used_in_production_version` - Whether the specified include is active on the production network and is referenced in parent's `production_version`.
