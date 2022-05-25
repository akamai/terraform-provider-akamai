---
layout: "akamai"
page_title: "Akamai: akamai_iam_countries"
subcategory: "Identity and Access Management"
description: |-
 IAM Countries
---

# akamai_iam_countries

Use `akamai_iam_countries` to retrieve all the possible countries that Akamai supports. Use the values from this data source to add or update a user's country information.

## Example usage

Basic usage:

```hcl
data "akamai_iam_countries" "countries" {
}

output "supported_countries" {
  value = data.akamai_iam_countries.countries
}
```

## Argument reference

There are no arguments for this data source.

## Attributes reference

These attributes are returned:

* `countries` — A list of countries.

[API Reference](https://developer.akamai.com/api/core_features/identity_management_user_admin/v2.html#getadmincountries)
