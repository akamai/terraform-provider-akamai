---
layout: "akamai"
page_title: "Akamai: akamai_iam_notification_prods"
subcategory: "IAM"
description: |-
 IAM Notification Products
---

# akamai_iam_notification_prods

Use `akamai_iam_notification_prods` to list all products a user can subscribe to and receive notifications for on the account. The account is determined by the tokens in your API client.

## Example Usage

Basic usage:

```hcl
data "akamai_iam_notification_prods" "notification_prods" {
}

output "supported_notification_prods" {
  value = data.akamai_iam_notification_prods.notification_prods
}
```

## Argument Reference

There are no arguments for this data source.

## Attributes Reference

The following attributes are returned:

* `products` — Products a user can subscribe to and receive notifications for on the account

[API Reference](https://developer.akamai.com/api/core_features/identity_management_user_admin/v2.html#getadminnotificationproducts)