---
layout: "akamai"
page_title: "Akamai: Get Started with Identity and Access Management"
description: |-
  Get Started with Akamai Identity and Access Management using Terraform
---

# Get Started with Identity and Access Management

The Akamai Provider for Terraform provides you the ability to automate the creation, and management of users, user notifications, and user grants.  

To get more information about Identity and Access Management, see:

* [API documentation](https://developer.akamai.com/api/core_features/identity_management_user_admin/v2.htm)
* How-to Guides
    * [Official Documentation](https://learn.akamai.com/en-us/products/core_features/identity_management.html)

## Configure the Terraform Provider

Set up your .edgerc credential files as described in [Get Started with Akamai APIs](https://developer.akamai.com/api/getting-started), and include read-write permissions for the Property Manager API. 

1. Create a new folder called `terraform`
1. Inside the new folder, create a new file called `akamai.tf`.
1. Add the provider configuration to your `akamai.tf` file:

```hcl
provider "akamai" {
	edgerc = "~/.edgerc"
	config_section = "papi"
}
```

## Prerequisites

To create a user there is a single dependencies you must first meet:

* **Country**: The user's country


## Retrieving Supported Countries

You can fetch fetch the list of supported countries using the [`iam_akamai_countries` data source](../data-sources/iam_supported_countries.md). To fetch the default contract ID no attributes need to be set:

## Creating a User

The user is represented by an [`iam_akamai_user` resource](../resources/user.md). Add this new block to your `akamai.tf` file after the provider block.

To define the entire configuration, we start by opening the resource block and give it a name. In this case we’re going to use the name "example".

Once you have a valid country, your user should look like this:

```hcl
resource "iam_akamai_user" "example" {

}
```

## Initialize the User

Once you have your configuration complete, save the file. Then switch to the terminal to initialize Terraform using the command:

```bash
$ terraform init
```

This command will install the latest version of the Akamai provider, as well as any other providers necessary (such as the local provider). To update the Akamai provider version after a new release, simply run `terraform init` again.

## Test Your Configuration

To test your configuration, use `terraform plan`:

```bash
$ terraform plan
```

This command will make Terraform create a plan for the work it will do based on the configuration file. This will not actually make any changes and is safe to run as many times as you like.

## Apply Changes

To actually create our property, we need to instruct Terraform to apply the changes outlined in the plan. To do this, in the terminal, run the command:

```bash
$ terraform apply
```

Once this completes your property will have been created. You can verify this in [Akamai Control Center](https://control.akamai.com) or via the [Akamai CLI](https://developer.akamai.com/cli). However, the property configuration has not yet been activated, so let’s do that next!

