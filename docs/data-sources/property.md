---
layout: "akamai"
page_title: "Akamai: akamai_property"
subcategory: "Provisioning"
description: |-
 Property
---

# akamai_property

Use the `akamai_property` data source to query and list the property ID and rule tree based on the property name.

## Example usage

This example returns the property ID and rule tree based on the property name and optional version argument:


```hcl
data "akamai_property" "example" {
    name = "terraform-demo"
    version = "1"
}

output "my_property_ID" {
  value = data.akamai_property.example
}
```

## Argument reference

This data source supports these arguments:

* `name` - (Required) The property name.
* `version` - (Optional) The property version whose ID you want to list.

## Attributes reference

This data source returns these attributes:

* `property_ID` - A property's unique identifier, including the `prp_` prefix.
* `rules` - A JSON-encoded rule tree for a given property.
