---
layout: "akamai"
page_title: "Akamai: akamai_property_contracts"
subcategory: "Provisioning"
description: |-
 Property contracts
---

# akamai_property_contracts


Use the `akamai_property_contracts` data source to list contracts associated with the [EdgeGrid API client token](https://developer.akamai.com/getting-started/edgegrid) you're using. 

## Example Usage

Return contracts associated with the EdgeGrid API client token currently used for authentication:

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

This data source returns these attributes:

* `contracts` - A list of supported contracts, with the following properties:
  * `contract_id` - the contract's unique ID. If your ID doesn't include the `ctr_` prefix, the Akamai Provider appends it 
  to your entry for processing purposes.
  * `contract_type_name` - The type of contract, either `DIRECT_CUSTOMER`, `INDIRECT_CUSTOMER`, `PARENT_CUSTOMER`,
  `REFERRAL_PARTNER`, `TIER_1_RESELLER`, `VAR_CUSTOMER`, `VALUE_ADDED_RESELLER`, `PARTNER`, `PORTAL_PARTNER`,
  `STREAMING_RESELLER`, `AKAMAI_INTERNAL`, or `UNKNOWN`.

