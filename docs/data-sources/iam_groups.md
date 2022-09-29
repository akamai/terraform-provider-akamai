---
layout: akamai
subcategory: Identity and Access Management
---

# akamai_iam_groups

Use `akamai_iam_groups` to list all groups in which you have a scope of admin for the current account and contract type. The account and contract type are determined by the access tokens in your API client. Use the `group_id` from this data source to construct the `auth_grants_json` when creating or updating a user's auth grants.

## Example usage

Basic usage:

```hcl
data "akamai_iam_groups" "my-groups" {}

output "groups" {
  value = data.akamai_iam_groups.my-groups
}
```

## Argument reference

There are no arguments for this data source.

## Attributes reference

This data source returns this attribute:

* `groups` — A set of groups for the contract, including:
  * `group_id` - Unique identifier for each group.
  * `name` - The name you supply for the group.
  * `parent_group_id` - For nested groups, identifies the parent group to which the current group belongs.
  * `sub_groups` - Set of nested Group objects.

[API Reference](https://techdocs.akamai.com/iam-api/reference/get-groups)
