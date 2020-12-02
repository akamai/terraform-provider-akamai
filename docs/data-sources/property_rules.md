---
layout: "akamai"
page_title: "Akamai: akamai_property_rules"
subcategory: "Provisioning"
description: |-
 Property ruletree
---

# akamai_property_rules


Use the `akamai_property_rules` data source to query and retrieve the rule tree of an existing property version.  Lets you 
search across contracts and groups you have access to.


## Basic Usage

This example returns the rule tree for version 3 of a property based on the selected contract and group:

```hcl-terraform
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

## Argument Reference

This data source supports these arguments:

* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix. 
* `group_id` - (Required) A group's unique ID, including the `grp_` prefix.
* `property_id` - (Required) A property's unique ID, including the `prp_` prefix. 
* `version` - (Optional) The version to return. Returns the latest version by default.

## Attributes Reference

This data source returns these attributes:

* `rules` - A JSON-encoded rule tree for the property.
* `errors` - A list of validation errors for the rule tree object returned. For more information see [Errors](https://developer.akamai.com/api/core_features/property_manager/v1.html#errors) in the Property Manager API documentation.
