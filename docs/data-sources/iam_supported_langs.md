---
layout: akamai
subcategory: Identity and Access Management
---

# akamai_iam_supported_langs

Use `akamai_iam_supported_langs` to list all the possible languages Akamai supports. Use the values from this API to set the preferred language for a user. Users should see Control Center in the language you set for them. The default language is English.

## Example usage

Basic usage:

```hcl
data "akamai_iam_supported_langs" "supported_langs" {
}

output "supported_supported_langs" {
  value = data.akamai_iam_supported_langs.supported_langs
}
```

## Argument reference

There are no arguments for this data source.

## Attributes reference

These attributes are returned:

* `languages` — Languages supported by Akamai

[API Reference](https://techdocs.akamai.com/iam-api/reference/get-user-languages)
