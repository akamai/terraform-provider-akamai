---
layout: "akamai"
page_title: "Akamai: akamai_iam_contact_types"
subcategory: "IAM"
description: |-
 IAM Contact Types
---

# akamai_iam_contact_types

Use `akamai_iam_contact_types` to retrieve all the possible `contact_types` that Akamai supports. Use the values from this data source to add or update a user’s contact type.

## Example usage

Basic usage:

```hcl
data "akamai_iam_contact_types" "contact_types" {
}

output "supported_contact_types" {
  value = data.akamai_iam_contact_types.contact_types
}
```

## Argument reference

There are no arguments for this data source.

## Attributes reference

These attributes are returned:

* `contact_types` — A list of contact types.

[API Reference](https://developer.akamai.com/api/core_features/identity_management_user_admin/v2.html#getadmincontacttypes)
