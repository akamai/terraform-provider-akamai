---
layout: "akamai"
page_title: "Akamai: akamai_property"
subcategory: "Provisioning"
description: |-
 Search property
---

# akamai_property


Use the `akamai_property` data source to query and retrieve general information about and the rule tree of an 
existing property.  With this data source you can search across the contracts and groups you have access to.

## Example Usage

This example returns property information and the rule tree for any properties associated with the `my-example.com` hostname:

datasource-example.tf
```hcl-terraform
datasource "akamai_property" "my-example" {
    search_type = "hostname"
    search_value   = "my-example.com"
}

output "property_match" {
  value = data.akamai_property.my-example
}
```

## Argument Reference

This data source supports these arguments:

* `search_type` - (Required) The item to search on. You can choose one of the following values: `name` for the property name, `hostname`, or `edge_hostname`. An `edge_hostname` search only includes active properties.
* `search_value` - (Required) The literal value to search on.

## Attributes Reference

This data source returns this attribute:

* `json` - The response to the query returned from the [Property Manager API](https://developer.akamai.com/api/core_features/property_manager/v1.html). 
