---
layout: "akamai"
page_title: "Akamai: Get Started with Identity and Access Management"
description: |-
  Get Started with Akamai Identity and Access Management using Terraform
---

# Get Started with Identity and Access Management Module

The Akamai Provider for Terraform lets you automate the creation and management of users, user notifications, and user grants.

To get more information about Identity and Access Management, see:

* [API and developer docs](https://developer.akamai.com/api/core_features/identity_management_user_admin/v2.html)
* [Product docs and resources](https://learn.akamai.com/en-us/products/core_features/identity_management.html)

## Prerequisites

Before you can create a user, you need:

* The user's country
* To complete the tasks in the 
[Get Started with the Akamai Terraform Provider](../guides/get_started_provider.md) 
guide.

## Identity and Access Management workflow
To set up the Identity and Access Management module, you need to:

* [Retrieve supported countries](#retrieve-supported-countries)
* [Create a user](#create-a-user)
* [Initialize the user](#initialize-the user)

## Retrieve supported countries

To fetch a list of supported countries, use the [`iam_akamai_countries` data source](../data-sources/iam_supported_countries.md). Attributes aren't needed to fetch the default contract ID.

## Create a user

The [`iam_akamai_user` resource](../resources/user.md) represents the user.

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
  send_otp_email = true
  auth_grants_json = jsonencode([
    {
      roleId = 3
      groupId = 12345
    }
  ])
}
```

## Initialize the user

After your configuration completes, save the file, then switch to the terminal to initialize Terraform using this command:

```bash
$ terraform init
```

This command installs the latest version of the Akamai Provider and any other providers necessary, such as the local provider. To update the Akamai Provider version after a new release, simply run `terraform init` again.
