---
layout: "akamai"
page_title: "Akamai: akamai_property_contracts"
subcategory: "Provisioning"
description: |-
 Property contracts
---

# akamai_property_contracts


Use `akamai_property_contracts` data source to list contracts associated with an edgerc API token. 

## Example Usage

Return what contracts exist for the user:

datasource-example.tf
```hcl-terraform
datasource "akamai_property_contracts" "my-example" {
}

output "property_match" {
  value = data.akamai_property_contracts.my-example
}
```

## Argument Reference

No arguments are supported:

## Attributes Reference

The following are the return attributes:

* `json` — PAPIs response to the query.

Example PAPI response is as follows:
```json
{
    "accountId": "act_1-9ZYX87",
    "contracts": {
        "items": [
            {
                "contractId": "ctr_1-1ABC123",
                "contractTypeName": "DIRECT_CUSTOMER"
            }
        ]
    }
}
```