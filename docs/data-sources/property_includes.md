---
layout: akamai
subcategory: Property Provisioning
---

# akamai_property_includes (Beta)

Use the `akamai_property_includes` data source to get all includes available for the current contract and group. Includes are small, reusable, and configurable components for your properties.

## Basic usage

This example returns all includes for the specified contract and group:

```hcl
data "akamai_property_includes" "my_example" {
    contract_id = "ctr_1-AB123"
    group_id    = "grp_12345"
}

output "my_example" {
  value = data.akamai_property_includes.my_example
}
```

## Argument reference

This data source supports these arguments:

* `contract_id` - (Required) A contract's unique ID, including the optional `ctr_` prefix.
* `group_id` - (Required) A group's unique ID, including the optional `grp_` prefix.
* `parent_property` - (Optional) The property that references the includes you want to list.
 * `id` - (Required) The property's unique identifier.
 * `version` - (Required) The version of the activated parent property.
* `type` - (Optional) Specifies the type of the include, either `MICROSERVICES` or `COMMON_SETTINGS`. Use this field for filtering. `MICROSERVICES` allow different teams to work independently on different parts of a single site. `COMMON_SETTINGS` includes are useful for configurations that share a large number of settings, often managed by a central team.

## Attributes reference

This data source returns these attributes:

* `includes` -  The small, reusable, configurable components for your properties.
 * `latest_version` - Returns the most recent version of the include.
 * `staging_version` - The include version currently activated on the staging network.
 * `production_version` - The include version currently activated on the production network.
 * `id` - The include's unique identifier.
 * `name` - The descriptive name for the include.
 * `type` - Specifies the type of the include, either `MICROSERVICES` or `COMMON_SETTINGS`.
