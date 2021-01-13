---
layout: "akamai"
page_title: "Akamai: akamai_property_products"
subcategory: "Provisioning"
description: |-
 Property products
---

# akamai_property_products


Use the `akamai_property_products` data source to list the products included on a contract. 

## Example usage

This example returns products associated with the [EdgeGrid client token](https://developer.akamai.com/getting-started/edgegrid) for a given contract:

```hcl
data "akamai_property_products" "my-example" {
    contract_id = "ctr_1-AB123"
}

output "property_match" {
  value = data.akamai_property_products.my-example
}
```

## Argument reference

This data source supports this argument:

* `contract_id` - (Required) A contract's unique ID, including the `ctr_` prefix. 

## Attributes reference

This data source returns these attributes:

* `products` - A list of supported products for the contract, including:
  * `product_id` - The product's unique ID, including the `prd_` prefix.
  * `product_name` - A string containing the product name.
