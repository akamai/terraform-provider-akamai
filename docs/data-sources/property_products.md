---
layout: "akamai"
page_title: "Akamai: akamai_property_products"
subcategory: "Provisioning"
description: |-
 Property products
---

# akamai_property_products


Use `akamai_property_products` data source to list products associated with a contract. 

## Example Usage

Given a contract return what products exist for the user:

datasource-example.tf
```hcl-terraform
datasource "akamai_property_products" "my-example" {
    contract_id = "ctr_1-AB123"
}

output "property_match" {
  value = data.akamai_property_products.my-example
}
```

## Argument Reference

The following arguments are supported:

* `contract_id` — (Required) The Contract ID.  Can be provided with or without `ctr_` prefix.

## Attributes Reference

The following are the return attributes:

* `json` — PAPIs response to the query.

Example PAPI response is of the form that follows:
```json
{
    "accountId": "act_1-9ZYX87",
    "contractId": "ctr_1-ABC234",
    "products": {
        "items": [
            {
                "productName": "Alta",
                "productId": "prd_Alta"
            }
        ]
    }
}

```