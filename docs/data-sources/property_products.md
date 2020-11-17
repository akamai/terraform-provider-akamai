---
layout: "akamai"
page_title: "Akamai: akamai_property_products"
subcategory: "Provisioning"
description: |-
 Property products
---

# akamai_property_products


Use the `akamai_property_products` data source to list the products included on a contract. 

## Example Usage

This example returns products associated with the [EdgeGrid client token](https://developer.akamai.com/getting-started/edgegrid) for a given contract:

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

This data source supports this argument:

* `contract_id` - (Required) A contract's unique ID. If your ID doesn't include the `ctr_` prefix, the Akamai Provider appends it to your entry for processing purposes. 

## Attributes Reference

This data source returns these attributes:

* `products` - A list of supported products for the contract, including:
  * `product_id` - A string containing product's unique ID. All results will include the `prd_` prefix.
  * `product_name` - A string containing the product name.
