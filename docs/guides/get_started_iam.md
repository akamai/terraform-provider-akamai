---
layout: "akamai"
page_title: "Module: Identity and Access Management"
description: |-
  Identity and Access Management module for the Akamai Terraform Provider
---

# Identity and Access Management Module Guide

The Identity and Access Management module lets you automate the creation and management of users, groups, and roles.

To get more information about Identity and Access Management, see the [product documentation](https://techdocs.akamai.com/iam/docs).

## Prerequisites

Before you can create a user, you need to:

1. Complete the tasks in the 
[Get Started with the Akamai Provider](../guides/get_started_provider.md) 
guide.
1. Set up your API client for [Identity and Access Management](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/akamai_provider_auth).
1. Retrieve supported countries and timezones using the [`iam_akamai_countries`](../data-sources/iam_countries.md) and [`akamai_iam_timezones`](../data-sources/iam_timezones.md) data sources. 

## Identity and Access Management workflows

Use Identity and Access Management to manage access privileges and users. When combined, users, groups, and roles grant access to Akamai applications, services, and objects. 

~> For more information about these concepts, see [API concepts](https://techdocs.akamai.com/iam-user-admin/reference/api-concepts) in the API documentation. 

For Identity and Access Management, there are three objects to create:

* [Users](#create-users) 
* [Roles](#create-roles)
* [Groups](#create-groups)

## Create users
To set up users, you need to:

The [`iam_akamai_user` resource](../resources/iam_user.md) represents the user.

To define the entire configuration, open the resource block and give it a name. For this case, you're going to use the name "example".

Once you have a valid country, your user should look like this:

```hcl
resource "iam_akamai_user" "example" {
  first_name = "John"
  last_name = "Doe"
  email = "john.doe@mycompany.com"
  country = "USA"
  phone = "(123) 321-1234"
  enable_tfa = false
  auth_grants_json = jsonencode([
    {
      roleId = 3
      groupId = 12345
    }
  ])
}
```

## Create roles 

Use [`akamai_iam_roles` resource](../resources/iam_roles.md) to set up the roles. 

To see if there are existing roles, start with the [`akamai_iam_grantable_roles` data source](../data-sources/iam_grantable_roles.md).

## Create groups

Use [`akamai_iam_group` resource](../resources/iam_group.md) to create a group. 

To see if there are existing groups, start with the [`akamai_iam_group` data source](../data-sources/iam_group.md). 
