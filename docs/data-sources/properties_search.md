---
layout: "akamai"
page_title: "Akamai: akamai_properties_search"
subcategory: "Property Provisioning"
description: |-
 Search
---

# akamai_properties_search


Use the `akamai_properties_search` data source to retrieve the list of properties matching a specific hostname, edge hostname or property name based on the [EdgeGrid API client token](https://techdocs.akamai.com/developer/docs/authenticate-with-edgegrid) you're using.

## Example usage

Return properties associated with the EdgeGrid API client token currently used for authentication:


```hcl
datasource "akamai_properties_search" "example" {
  key = "hostname"
  value = "test.akamai.com"
}

output "my_property_list" {
  value = data.akamai_properties_search.example
}
```

## Argument reference

This data source supports these arguments:

* `key` - (Required) Key used for search. Valid values are:
  * **hostname**
  * **edgeHostname**
  * **propertyName**
* `value` - (Required) Value to search for.

## Attributes reference

This data source returns this attribute:

* `properties` - A list of property version matching the given criteria.