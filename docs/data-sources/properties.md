---
layout: "akamai"
page_title: "Akamai: akamai_properties"
subcategory: "Provisioning"
description: |-
 Properties
---

# akamai_properties


Use `akamai_properties` data source to query and retrieve the list of properties for a group and contract 
that the current API token has access to. 

## Example Usage

Given a contract and group return what properties exist for the user:

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

The following are the return attributes:

* `json` — PAPIs response to the query.
