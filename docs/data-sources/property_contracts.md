---
layout: "akamai"
page_title: "Akamai: akamai_property_contracts"
subcategory: "Common"
description: |-
 Property contracts
---

# akamai_property_contracts


Use the `akamai_property_contracts` data source to list contracts associated with the [EdgeGrid API client token](https://developer.akamai.com/getting-started/edgegrid) you're using. 

## Example usage

Return contracts associated with the EdgeGrid API client token currently used for authentication:

```hcl
data "akamai_property_contracts" "my-example" {
}

output "property_match" {
  value = data.akamai_property_contracts.my-example
}
```

## Argument reference

There are no arguments available for this data source.

## Attributes reference

This data source returns these attributes:

* `contracts` - A list of supported contracts, with the following properties:
  * `contract_id` - The contract's unique ID, including the `ctr_` prefix.
  * `contract_type_name` - The type of contract, either `DIRECT_CUSTOMER`, `INDIRECT_CUSTOMER`, `PARENT_CUSTOMER`, `REFERRAL_PARTNER`, `TIER_1_RESELLER`, `VAR_CUSTOMER`, `VALUE_ADDED_RESELLER`, `PARTNER`, `PORTAL_PARTNER`, `STREAMING_RESELLER`, `AKAMAI_INTERNAL`, or `UNKNOWN`.

