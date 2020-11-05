---
layout: "akamai"
page_title: "Akamai: akamai_property"
subcategory: "Provisioning"
description: |-
 Search property
---

# akamai_property


Use `akamai_property` data source to query and retrieve the instance information and rule tree of an 
existing property instance.  allows searching across contracts and groups you may have access to.

## Example Usage

Given a contract and group return what properties exist for the user:

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

The following arguments are supported:

* `search_type` — (Required) One of the following values `name`, `hostname`, or `edge_hostname` (last field only searches active properties).
* `search_value` — (Required) The value to be searched for in the field specified.

## Attributes Reference

The following are the return attributes:

* `json` — PAPIs response to the query.
