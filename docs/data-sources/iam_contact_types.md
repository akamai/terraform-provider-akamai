---
layout: "akamai"
page_title: "Akamai: akamai_iam_contact_types"
subcategory: "IAM"
description: |-
 IAM Contact Types
---

# akamai_iam_contact_types

Use `akamai_iam_contact_types` datasource to retrieve all the possible contact_types that Akamai supports. Use the values from this operation to add or update a user’s contactType.

## Example Usage

Basic usage:

```hcl
data "akamai_iam_contact_types" "contact_types" {
}

output "supported_contact_types" {
  value = data.akamai_iam_contact_types.contact_types
}
```

## Argument Reference

There are no arguments for this data source.

## Attributes Reference

The following attributes are returned:

* `contact_types` — A list of contact types

[API Reference](https://developer.akamai.com/api/core_features/identity_management_user_admin/v2.html#getadmincontacttypes)