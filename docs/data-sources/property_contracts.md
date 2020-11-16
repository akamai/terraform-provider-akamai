---
layout: "akamai"
page_title: "Akamai: akamai_property_contracts"
subcategory: "Provisioning"
description: |-
 Property contracts
---

# akamai_property_contracts


Use `akamai_property_contracts` data source to list contracts associated with an EdgeGrid API client token. 

## Example Usage

Return contracts associated with the EdgeGrid API client token:

datasource-example.tf
```hcl-terraform
datasource "akamai_property_contracts" "my-example" {
}

output "property_match" {
  value = data.akamai_property_contracts.my-example
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

The following attributes are returned:

* `contracts` — list of supported contracts, with the following properties:
  * `contract_id` - the contract ID (string)
  * `contract_type_name` - the contract type (string)
