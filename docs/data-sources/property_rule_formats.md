---
layout: "akamai"
page_title: "Akamai: akamai_property_rule_formats"
subcategory: "Provisioning"
description: |-
 Properties rule formats
---

# akamai_property_rule_formats

Use the `akamai_property_rule_formats` data source to query the list of known rule formats. 
You use rule formats to [freeze](https://developer.akamai.com/api/core_features/property_manager/v1.html#freezerf) or 
[update](https://developer.akamai.com/api/core_features/property_manager/v1.html#updaterf) the versioned set of behaviors
and criteria a rule tree invokes. Without this mechanism, behaviors and criteria would update automatically and 
generate unexpected errors. 


## Example usage

Use this example to list available property rule formats:

```hcl-terraform
datasource "akamai_property_rule_formats" "my-example" {
}

output "property_match" {
  value = data.akamai_property_rule_formats.my-example
}
```

## Argument reference

There are no arguments available for this data source.

## Attributes reference

This data source returns this attribute:

* `formats` - A list of supported rule format identifiers. For example: 

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
