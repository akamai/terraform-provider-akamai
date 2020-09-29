---
layout: "akamai"
page_title: "Akamai: Get Started with Property Management"
subcategory: "docs-akamai-guide-get-started-property"
description: |-
  Get Started with Akamai Property Management using Terraform
---

# Get Started with Property Management

The Akamai Provider for Terraform provides you the ability to automate the creation, deployment, and management of property configuration and activation, edge hostnames, and CP Codes.  

To get more information about Property Management, see:

* [API documentation](https://developer.akamai.com/api/core_features/property_manager/v1.html)
* How-to Guides
    * [Official Documentation](https://learn.akamai.com/en-us/products/core_features/property_manager.html)

## Configure the Terraform Provider

Set up your .edgerc credential files as described in [Get Started with Akamai APIs](https://developer.akamai.com/api/getting-started), and include read-write permissions for the Property Manager API. 

1. Create a new folder called `terraform`
1. Inside the new folder, create a new file called `akamai.tf`.
1. Add the provider configuration to your `akamai.tf` file:

```hcl
provider "akamai" {
    edgerc = "~/.edgerc"
    papi_section = "papi"
}
```

## Prerequisites

To create a property there are a number of dependencies you must first meet:

* **Contract ID**: The ID of the contract under which the property, CP Code, and edge hostnames will live
* **Group ID**: The ID of the group under which the property, CP Code, and edge hostnames will live
* **Edge hostname:** The Akamai edge hostname for your property. You can [create a new one or reuse an existing one](#managing-edge-hostnames). 
* **Origin hostname:** The origin hostname you want your configuration to point to
* **Product:** The [Akamai Product ID](/docs/providers/akamai/g/appendix.html#common-product-ids) for the product you are using (Ion, DSA, etc.)
* **Rules configuration**: The rules.json file contains the base rules for the property.  (learn how to leverage the rules.json tree from an existing property [here](/docs/providers/akamai/g/faq.html#migrating-a-property-to-terraform))


## Retrieving The Contract ID

You can fetch your contract ID automatically using the [`akamai_contract` data source](/docs/providers/akamai/d/contract.html). To fetch the default contract ID no attributes need to be set:

```hcl
data "akamai_contract" "default" {

}
```

Alternatively, if you have multiple contracts, you can specify the `group` which contains it:

```hcl
data "akamai_contract" "default" {
  group = "default"
}
```

You can now refer to the contract ID using the `id` attribute: `data.akamai_contract.default.id`.

## Retrieving The Group ID

Similarly, you can fetch your group ID automatically using the [`akamai_group` data source](/docs/providers/akamai/d/group.html). To fetch the default group ID no attributes need to be set:

```hcl
data "akamai_group" "default" {

}
``` 

To fetch a specific group, you can specify the `name` argument:

```hcl
data "akamai_group" "default" {
  name = "example"
}
```

You can now refer to the group ID using the `id` attribute: `data.akamai_group.default.id`.

## Managing Edge Hostnames

Whether you are reusing an existing Edge Hostname or creating a new one, you use the `akamai_edge_hostname` resource. The following will create the `example.com.edgesuite.net` edge hostname:

```hcl
resource "akamai_edge_hostname" "example" {
    group = "${data.akamai_group.default.id}"
	contract = "${data.akamai_contract.default.id}"
	product = "prd_SPM"
    edge_hostname = "example.com.edgesuite.net"
}
```

> **Note:** Notice that we’re using variables from the previous section to reference the group and contract IDs. These will automatically be replaced at runtime by Terraform with the actual values.

This will create a non-secure hostname, to create a secure hostname, you must specify a certificate enrollment ID, using the `certificate` argument:

```hcl
resource "akamai_edge_hostname" "example" {
    group = "${data.akamai_group.default.id}"
	contract = "${data.akamai_contract.default.id}"
	product = "prd_SPM"
    edge_hostname = "example.com.edgesuite.net"
    certificate = "<CERTIFICATE ENROLLMENT ID>"
}
```

This will create a Standard TLS secure hostname, to create an Enhanced TLS hostname, use the `edgekey.net` domain suffix for the `edge_hostname` instead.

> **Note:** This resource does not automatically make the property secure. You will need the `is_secure` flag set to `true` in your rule tree as well — this can be set in your `akamai_property` or `akamai_property_rules` resources, or in your `rules.json` file.

## Property Rules

A property contains the delivery configuration, or rule tree, which determines how requests are handled. This rule tree is usually represented using JSON, and is often refered to as `rules.json`.

You can specify the rule tree as a JSON string, using the [`rules` argument of the `akamai_property` resource](/docs/providers/akamai/r/property.html#rules).

We recommend storing the rules JSON as a JSON file on disk and ingesting it using Terraforms `local_file` data source. For example, if our file is called `rules.json`, we might create a `local_file` data source called `rules`. We specify the path to `rules.json` using the `filename` argument:

```hcl
data "local_file" "rules" {
    filename = "rules.json"
}
```

We can now use `${data.local_file.rules.content}` to reference the file contents in the `akamai_property.rules` argument.

## Creating a Property

The property itself is represented by an [`akamai_property` resource](/docs/providers/akamai/r/property.html). Add this new block to your `akamai.tf` file after the provider block.

To define the entire configuration, we start by opening the resource block and give it a name. In this case we’re going to use the name "example".

Next, we set the name of the property, contact email id, product ID, group ID, CP code, property hostname, and edge hostnames.


Finally, we setup the property rules: first, we should specify the [`rule format` argument](/docs/providers/akamai/r/property.html#rule_format), as well as passing the `rules.json` data to `rules` argument.

Once you’re done, your property should look like this:

```hcl
resource "akamai_property" "example" {
	name = "xyz.example.com"                        # Property Name
	contact = ["user@example.org"]                  # User to notify of de/activations  
	product  = "prd_SPM"                            # Product Identifier (Ion)
	group    = "${data.akamai_group.default.id}"    # Group ID variable
	contract = "${data.akamai_contract.default.id}" # Contract ID variable
	hostnames = {                                   # Hostname configuration
	    # "public hostname" = "edge hostname"
        "example.com" = "example.com.edgesuite.net"
        "www.example.com" = "example.com.edgesuite.net"
    }
	rule_format = "v2018-02-27"                     # Rule Format
	rules = "${data.local_file.rules.content}"      # JSON Rule tree
}
```

> **Note:** If you are creating a secure property (using TLS), you need to set the `is_secure` attribute to true unless it already set in your `rules.json`. If specified in the property resource, it will *override* the value in `rules.json`.


## Initialize the Provider

Once you have your configuration complete, save the file. Then switch to the terminal to initialize terraform using the command:

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

To actually create our property, we need to instruct terraform to apply the changes outlined in the plan. To do this, in the terminal, run the command:

```bash
$ terraform apply
```

Once this completes your property will have been created. You can verify this in [Akamai Control Center](https://control.akamai.com) or via the [Akamai CLI](https://developer.akamai.com/cli). However, the property configuration has not yet been activated, so let’s do that next!

## Activate your property


To activate your property we need to create a new [`akamai_property_activation` resource](/docs/providers/akamai/r/property_activation.html). This resource manages the activation for a property, allowing you to specify which network and what version to activate.

You will need to set the `property` ID and `version` arguments, which can both be set from the `akamai_property` resource. You should then set the `network` to `STAGING` or `PRODUCTION`. You should also set the `contact` email address.

Lastly, you need to affirm that you wish to activate the property, by setting the `activate` argument to `true`.

```hcl
resource "akamai_property_activation" "example" {
	property = "${akamai_property.example.id}"
	version = "${akamai_property.example.version}"
	network = "STAGING"
	contact = ["user@example.org"]
	activate = true	
}
```

### Test & Deploy Property Activation

Like the property itself, we should test our configuration with this command:

```bash
$ terraform plan
```

This time you will notice how the property is not being modified while the activation is being added to the plan.

Again, as with our property configuration we can apply our changes using:

```bash
$ terraform apply
```

This will activate the property on the staging network.
