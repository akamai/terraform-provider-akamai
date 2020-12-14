---
layout: "akamai"
page_title: "Akamai: akamai_iam_groups"
subcategory: "IAM"
description: |-
 IAM Groups
---

# akamai_iam_groups

Use `akamai_iam_groups` to list all groups in which you have a scope of admin for the current account and contract type. The account and contract type are determined by the access tokens in your API client.

## Example Usage

Basic usage:

```hcl
data "akamai_iam_groups" "my-groups" {
    group_id = "1234567"
    actions = true
}

output "groups" {
  value = data.akamai_iam_groups.my-groups
}
```

## Argument Reference

The following arguments are supported:

* `actions` - (optional, bool) When enabled, the response includes information about actions such as "edit" or "delete"

## Attributes Reference

The following attributes are returned:

* `groups` — A list of groups

[API Reference](https://developer.akamai.com/api/core_features/identity_management_user_admin/v2.html#getgroups)