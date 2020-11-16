---
layout: "akamai"
page_title: "Akamai: akamai_properties"
subcategory: "Provisioning"
description: |-
 Properties
---

# akamai_properties


Use `akamai_properties` data source to query and retrieve the list of properties for a group and contract 
that the current EdgeGrid API client token has access to. 

## Example Usage

Return properties associated with the EdgeGrid API client token given a contract and group:

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

The following arguments are supported:

* `contract_id` — (Required) The Contract ID.  Can be provided with or without `ctr_` prefix.
* `group_id` — (Required) The Group ID. Can be provided with or without `grp_` prefix.

## Attributes Reference

The following attributes are returned:

* `properties` — Provisioning API response containing a list of properties available for the provided contract and group.
