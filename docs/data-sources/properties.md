---
layout: "akamai"
page_title: "Akamai: akamai_properties"
subcategory: "Property Provisioning"
description: |-
 Properties
---

# akamai_properties


Use the `akamai_properties` data source to query and retrieve the list of properties for a group and contract
based on the [EdgeGrid API client token](https://techdocs.akamai.com/developer/docs/authenticate-with-edgegrid) you're using.

## Example usage

Return properties associated with the EdgeGrid API client token currently used for authentication:


```hcl
datasource "akamai_properties" "example" {
    contract_id = "ctr_1-AB123"
    group_id   = "grp_12345"
}

output "my_property_list" {
  value = data.akamai_properties.example
}
```

## Argument reference

This data source supports these arguments:

* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix.
* `group_id` - (Required) A group's unique ID, including the `grp_` prefix.

## Attributes reference

This data source returns this attribute:

* `properties` - A list of properties available for the contract and group IDs provided.
