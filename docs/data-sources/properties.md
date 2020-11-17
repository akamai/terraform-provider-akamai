---
layout: "akamai"
page_title: "Akamai: akamai_properties"
subcategory: "Provisioning"
description: |-
 Properties
---

# akamai_properties


Use the `akamai_properties` data source to query and retrieve the list of properties for a group and contract 
based on the [EdgeGrid API client token](https://developer.akamai.com/getting-started/edgegrid) you're using. 

## Example Usage

Return properties associated with the EdgeGrid API client token currently used for authentication:


datasource-example.tf
```hcl-terraform
datasource "akamai_properties" "example" {
    contract_id = "ctr_1-AB123"
    group_id   = "grp_12345"
}

output "my_property_list" {
  value = data.akamai_properties.example
}
```

## Argument Reference

This data source supports these arguments:

* `contract_id` - (Required) A contract's unique ID. If your ID doesn't include the `ctr_` prefix, the Akamai Provider appends it to your entry for processing purposes. 
* `group_id` - (Required) A group's unique ID. If your ID doesn't include the `grp_` prefix, the Akamai provider appends it to your entry for processing purposes.

## Attributes Reference

This data source returns these attributes:

* `properties` - A list of properties available for the contract and group IDs provided.
