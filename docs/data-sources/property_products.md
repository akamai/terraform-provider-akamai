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

Return products associated with the EdgeGrid API client token under a given contract:

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

The following attributes are returned:

* `products` — list of supported product, with the following properties:
  * `product_id` - the product ID (string)
  * `product_name` - the product name (string)
