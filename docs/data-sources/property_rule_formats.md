---
layout: "akamai"
page_title: "Akamai: akamai_property_rule_formats"
subcategory: "Provisioning"
description: |-
 Properties rule formats
---

# akamai_property_rule_formats


Use `akamai_property_rule_formats` data source to query the list of known rule formats.  These formats can be used to 'freeze'
the API syntax at a known format and avoid syntax changes that the `latest` format will undergo where behaviors and 
criteria undergo structural changes.

## Example Usage

List current property rule formats:

datasource-example.tf
```hcl-terraform
datasource "akamai_property_rule_formats" "my-example" {
}

output "property_match" {
  value = data.akamai_property_rule_formats.my-example
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

The following attributes are returned:

* `formats` — list of supported format identifiers.

Example response is of the form that follows:
```json
[
    "latest",
    "v2015-08-17",
    "v2015–08–17",
    "v2016–11–15",
    "v2017–06–19",
    "v2018–02–27",
    "v2018–09–12",
    "v2019–07–25",
    "v2020–03–04",
    "v2020–11–01"
]
```
