---
layout: "akamai"
page_title: "Akamai: akamai_property_rules"
subcategory: "Provisioning"
description: |-
 Property rule tree
---

# akamai_property_rules

~> **Note** Version 1.0.0 of the Akamai Terraform Provider is now available for the Property Provisioning module. To upgrade to the new version, you have to update this data source. See [Upgrade to Version 1.0.0](../guides/1.0_migration.md) for details. 

Use the `akamai_property_rules` data source to query and retrieve the rule tree of 
an existing property version. This data source lets you search across the contracts 
and groups you have access to.

## Basic usage

This example returns the rule tree for version 3 of a property based on the selected contract and group:

```hcl
datasource "akamai_property_rules" "my-example" {
    property_id = "prp_123"
    group_id = "grp_12345"
    contract_id = "ctr_1-AB123"
    version   = 3
}

output "property_match" {
  value = data.akamai_property_rules.my-example
}
```

## Argument reference

This data source supports these arguments:

* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix. 
* `group_id` - (Required) A group's unique ID, including the `grp_` prefix.
* `property_id` - (Required) A property's unique ID, including the `prp_` prefix. 
* `version` - (Optional) The version to return. Returns the latest version by default.

## Attributes reference

This data source returns these attributes:

* `rules` - A JSON-encoded rule tree for the property.
* `errors` - A list of validation errors for the rule tree object returned. For more information see [Errors](https://developer.akamai.com/api/core_features/property_manager/v1.html#errors) in the Property Manager API documentation.
