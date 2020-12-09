---
layout: "akamai"
page_title: "Akamai: Get Started with Property Management"
description: |-
  Get Started with Akamai Property Management using Terraform
---

# Get Started with the Provisioning Module

You can use Provisioning module resources and data sources to create,
deploy, activate, and manage properties, edge hostnames, and content
provider (CP) codes.

For more information about properties, see [Property Manager documentation](https://learn.akamai.com/en-us/products/core_features/property_manager.html).

## Prerequisites

Before you can create a property, you need to complete the tasks in the [Get Started with the Akamai Terraform Provider](../guides/get_started_provider.md) guide. Be sure you have the contract and group IDs you retrieved available. You'll need them to set up the Provisioning module.

## Provisioning Workflow 

To set up the Provisioning module, you need to: 

* [Retrieve the product ID](#retrieve-the-product-id). This is the ID for the product you are using, like Ion or Adaptive Media Delivery.
* [Add or create an edge hostname](#add-an-edge-hostname).
* [Set up rules for your property](#set-up-property-rules). A separate `rules.json` file contains the base rules for the property. 
* [Import or create a property](#import-or-create-a-property).
* [Apply your property changes](#apply_your_property_changes). This step adds the property to your Terraform configuration.
* [Activate your property](#activate_your_property]). Once you apply your property changes, you have to activate the property configuration for it to be live.

## Retrieve the product ID

When setting up properties, you need to retrieve the ID for the specific
Akamai product you are using. See the [Akamai Product ID](../guides/appendix.md#common-product-ids) section for a list of common IDs.

-> **Note** If you're currently using prefixes with your IDs, you might have to remove the `prd_` prefix from your entry. For more information about prefixes, see the [ID prefixes](https://developer.akamai.com/api/core_features/property_manager/v1.html#prefixes) section of the Property Manager API (PAPI) documentation.

## Add an edge hostname

You use the [akamai_edge_hostname](../resources/property_edge_hostname.md) resource to 
reuse an existing edge hostname or create a new one. 

To create different hostname types, you need to change the domain suffix
for the `edge_hostname` attribute. See [Domain Suffixes for Different Edge Hostname Types](../guides/appendix.md#domain-suffixes-for-different-edge-hostname-types)

Once you set up the `akamai_edge_hostname` resource, run `terraform plan` and resolve any errors or warnings before continuing with the next step. See [Command: plan](https://www.terraform.io/docs/commands/plan.html) for more information about this Terraform command.


### Non-secure hostnames

The following code will create a non-secure Standard TLS edge hostname for
`example.com`:

```hcl
resource "akamai_edge_hostname" "example" {
	group_id = data.akamai_group.default.id
	contract_id = data.akamai_contract.default.id
	product_id = "prd_SPM"
	edge_hostname = "example.com.edgesuite.net"
}
```

-> **Note** In this example, we're using variables for `contract_id` and `group_id`
 to reference the default group and contract. At runtime, Terraform automatically
 replaces these variables with the actual values.

### Secure hostnames

To create a secure hostname, you also need the `certificate` attribute, which is 
the certificate enrollment ID you can retrieve from the [Certificate Provisioning System CLI](https://github.com/akamai/cli-cps). Here's an example

```hcl
resource "akamai_edge_hostname" "example" {
	group_id = data.akamai_group.default.id
	contract_id = data.akamai_contract.default.id
	product_id = "prd_SPM"
	edge_hostname = "example.com.edgesuite.net"
	certificate = "1000"
}
```

As a final step, you'll have to set the `is_secure` flag to `true` in the `akamai_property_rules` resource. If you set this flag it *overrides* the value in the `rules.json` file.
<!--Is the paragraph above accurate? I didn't see this flag in the resource description.-->

## Set up property rules

A property contains the delivery configuration, or rule tree, that determines how Akamai handles requests. This rule tree is usually represented using JSON, and is often refered to as `rules.json`.

You can specify the rule tree as a JSON string, using the [`rules` argument of the `akamai_property` resource](../resources/property.md#rules).

As a best practice, you should store the rule tree as a JSON file on disk and ingest it using Terraform's `file` function. For example, if you name your file `rules.json`, you set up the `file` function like this:

```hcl
locals {
	json = file("${path.module}/rules.json") 
}
```

You can now use `local.json` to reference the file contents in the `akamai_property.rules` argument. Or you can embed the file reference directly in the `akamai_property` resource using the `rules` attribute: `rules = file("${path.module}/rules.json")`.

Before continuing with the next step, run `terraform plan` and resolve any errors or warnings. See [Command: plan](https://www.terraform.io/docs/commands/plan.html) for more information about this Terraform command.

## Import or create a property
You can either import an existing property or create a new one with Terraform: 

### Import a property

To import an existing property into Terraform you have to export the `rules.json`
 file from the property you want to import. You can use PAPI, the [Property Manager
 CLI](https://github.com/akamai/cli-property-manager), or the Property Manager application in Control Center to export this JSON file.

You'll then need to create an `akamai_property` resource that pulls in the `rules.json`.

You can use the `akamai_property_rules` data source generate a rule template. It reads the server's copy of the rules then generates output in a format that you can save in a JSON file. If your rule template includes variables, you'll have to set them up again.
<!--How's this new paragraph? Should we send them to the section in the migration guide that talks about variables? I'm pretty sure that guide has the most up-to-date info right now.-->

### Create a property

You use the [akamai_property
resource](../resources/property.md)
to represent your property. Add this new block to your `akamai.tf` file
after the `provider` block.

To define the entire configuration, start by opening the resource block
and give it a name, like `example`. Within the new block, you set the name 
of the property, contact email, product ID, group ID, CP code, property hostname,
 and edge hostnames.

Finally, you set up the property rules. You first specify the [rule
format argument](../resources/property.md#rule_format),
then add the path to the `rules.json` file. You can set a variable for the path, like `${path.module}`. 

Once you're done, your property should look like this:

```hcl
resource "akamai_property" "example" {
	name = "xyz.example.com"                        # Property Name
	contact = ["user@example.org"]                  # User to notify of de/activations  
	product_id  = "prd_SPM"                         # Product Identifier (Ion)
	group_id    = data.akamai_group.default.id      # Group ID variable
	contract_id = data.akamai_contract.default.id   # Contract ID variable
	hostnames = {                                   # Hostname configuration
		# "public hostname" = "edge hostname"
		"example.com" = "example.com.edgesuite.net"
		"www.example.com" = "example.com.edgesuite.net"
	}
	rule_format = "v2018-02-27"                     # Rule Format
	rules = file("${path.module}/rules.json")       # JSON Rule tree
}
```

Before continuing with the next step, run `terraform plan` and resolve any errors or warnings. See [Command: plan](https://www.terraform.io/docs/commands/plan.html) for more information about this Terraform command.
 

## Apply your property changes

To actually add the property to your Terraform configuration, you need to 
tell Terraform to apply the changes outlined in the plan. To do this, run
 `terraform apply` in the terminal.

Once the command completes your new property is available. You can verify
this in [Akamai Control Center](https://control.akamai.com/) or by using the
[Property Manager CLI](https://github.com/akamai/cli-property-manager). 

However, you still have to activate the property configuration for it to be live.

## Activate your property

Once you’re satisfied with any version of a property, an activation deploys it, 
either to the Akamai staging or production network. You activate a specific version, 
but the same version can be activated separately more than once. You can either 
cancel an activation shortly after requesting it, or in many cases, use a fast 
fallback feature within a matter of seconds to roll back a live activation to the 
previous activated version.
<!--Does the last sentence above apply to Terraform configs?-->

### Create your property activation resource

To activate your property you need to create a new
[akamai_property_activation
resource](../resources/property_activation.md).
This resource manages property activations, letting you specify the
property version to activate and the network to activate it on.

You need following to set up this resource:

* property ID and version, which you can set from the `akamai_property` resource.

* network, which is either `STAGING` or `PRODUCTION`. You should activate the property on staging first to verify that everything works as expected before activating on production. 

* The email addresses to send activation updates to.

Here's an example:

```hcl
resource "akamai_property_activation" "example" {
	property_id = akamai_property.example.id
	version = akamai_property.example.version
	network = "STAGING"
	contact = ["user@example.org"]
}
```

### Test and deploy your property activation

Like you did with the `akamai_property` resource, you should first verify the
[akamai_property_activation](https://registry.terraform.io/docs/providers/akamai/r/property_activation)
resource by running the `terraform plan` command:

The plan command adds the activation to the property. It doesn't change any other property settings.

If everything looks good, run `terraform apply` to start the activation. This command will activate the property on the selected network.

## How you can use Provisioning resources

The sections that follow include information on different ways to use
Provisioning resources.

## Dynamic Rule Trees Using Templates

You can use rule templates to insert values from Akamai Provider resources and data sources into your `rules.json` file.

The `akamai_property_rules_template` data source supports variable 
replacement and the use of snippet templates from the Property Manager CLI :

You may need to use different rule sets with different properties. To do this you 
need to maintain a base rule set and then import individual rule sets. You'll first need to create a directory structure, for example:

```dir
rules/main.json
rules/snippets/routing.json
rules/snippets/performance.json
…
```

The `rules` directory contains a single file, `main.json` and a `snippets` 
subdirectory that contain all of the smaller JSON rule files, or snippets. In the 
`main.json` file, you set up a basic template for our JSON, like this:

```json
{
    rules": {
      "name": "default",
      "children": [
        "#include:snippets/performance.json",
        "#include:snippets/routing.json"
      ],
      "options": {
            "is_secure": “${env.secure}"
      }
    },
    "ruleFormat": "v2018-02-27"
}
```

This enables our rules template to process the `rules.json` and pull each 
fragment that's referenced, like `routing.json`. The rendered output can be set 
into `akamai_property` resources.

```hcl
data "akamai_property_rules_template" "example" {
  template_file = abspath("${path.root}/rules/main.json")
  variables {
      name  = "secure"
      value = "true"
      type  = "bool"
  }
  variables {
      name  = "caching_ttl"
      value = "3d"
      type  = "string"
  }
}

resource "akamai_property" "example" {
    ....
    rules  = data.akamai_property_rules_template.example.json
}
```
