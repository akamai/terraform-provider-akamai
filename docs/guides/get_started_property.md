---
layout: "akamai"
page_title: "Akamai: Get Started with Property Management"
description: |-
  Get Started with Akamai Property Management using Terraform
---

# Get Started with Property Management

You can use Provisioning resources and data sources to create,
deploy, activate, and manage properties, edge hostnames, and content
provider codes (CP codes).

For more information about properties, see the [Property Manager documentation](https://learn.akamai.com/en-us/products/core_features/property_manager.html) page.

## Prerequisites

To create a property there are a number of dependencies you must first
meet:

* **Complete the tasks in Get Started.** You need to complete the tasks in the [Get    Started](https://registry.terraform.io/docs/providers/akamai/g/get_started_provider) guide. Be sure you have the contract ID and group ID you retrieved available. You'll need them to set up the Provisioning module.
	<!--Did I get the URL right?-->
* **Retrieve the Product ID**. You'll need the [Akamai Product ID](https://registry.terraform.io/docs/providers/akamai/g/appendix#common-product-ids) for the product you are using, like Ion or Adaptive Media Delivery.

<!--Go back to this. Need prereqs and a workflow.-->

* **Edge hostname:** The Akamai edge hostname for your property. You can [create a new one or reuse an existing one](#add-an-edge-hostname). 
<!--Did I get the URL right?-->
* **Origin hostname:** The origin hostname you want your property to point to. Your property should point to an origin hostname you create.
<!--Where do you use this? It doesn't seem to be mentioned later.-->
* **Rules configuration**: The `rules.json` file contains the base rules for the property. 

## Retrieve the product ID

When setting up properties, you need to retrieve the ID for the specific
Akamai product you are using. See the [Akamai Product ID](https://registry.terraform.io/docs/providers/akamai/g/appendix#common-product-ids) section for a list of common IDs.

-> **Note** If you're currently using prefixes with your IDs, you might have to remove the `prd_` prefix from your entry. For more information about prefixes, see the [ID prefixes](https://developer.akamai.com/api/core_features/property_manager/v1.html#prefixes) section of the PAPI documentation.

## Add an edge hostname

You use the [akamai_edge_hostname](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/edge_hostname) resource to reuse an existing edge
hostname or create a new one. 

To create different hostname types, you need to change the domain suffix
for the `edge_hostname` attribute. See [Domain Suffixes for Different Edge Hostname Types](https://registry.terraform.io/docs/providers/akamai/g/appendix#domain-suffixes-for-different-edge-hostname-types)

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

-> **Note** Notice that for the contract and group IDs we're using
variables to reference the default group and contract. At runtime,
Terraform automatically replaces these variables with the actual
values.

### Secure hostnames

To create a secure hostname, you also need to add the `certificate` attribute, which is the certificate enrollment ID from the [Certificate
    Provisioning System CLI](https://github.com/akamai/cli-cps). Here's an example

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

## Set up property rules

A property contains the delivery configuration, or rule tree, that determines how Akamai handles requests. This rule tree is usually represented using JSON, and is often refered to as `rules.json`.

You can specify the rule tree as a JSON string, using the [`rules` argument of the `akamai_property` resource](../resources/property.md#rules).

We recommend storing the rule tree as a JSON file on disk and ingesting it using Terraform's `file` function. For example, if you name your file `rules.json`, you set up the `file` function like this:

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
<!--Not sure this section is accurate. Taken from the current Migrate a Property section--> 

To import an existing property into Terraform you have to export the `rules.json` file from the property you want to import. You can use PAPI, the Property Manager CLI, or the Property Manager application Control Center to export this JSON file.

You'll then need to create an `akamai_property` resource that pulls in the `rules.json`.
<!--How would you do this exactly? I took a guess about the resource here.-->


### Create a property

You use the [akamai_property
resource](https://registry.terraform.io/docs/providers/akamai/r/property)
to represent your property. Add this new block to your `akamai.tf` file
after the `provider` block.

To define the entire configuration, start by opening the resource block
and give it a name, like `example`. Within the new block, you set the name of the property, contact email, product ID, group
ID, content provider (CP) code, property hostname, and edge hostnames.

Finally, you set up the property rules. You first specify the [rule
format argument](https://registry.terraform.io/docs/providers/akamai/r/property#rule_format),
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

-> **Note** If you're creating a secure property using Enhanced TLS, you'll have to set the `is_secure` flag to `true` in the `akamai_property` resource. If you set this flag it *overrides* the value in the `rules.json` file.
<!--Is this the correct resource?-->

Before continuing with the next step, run `terraform plan` and resolve any errors or warnings. See [Command: plan](https://www.terraform.io/docs/commands/plan.html) for more information about this Terraform command.
 

## Apply your property changes
<!--Do we have to run `terraform init` at all?-->

To actually add the property to your Terraform configuration, we need to instruct Terraform to apply
the changes outlined in the plan. To do this, run `terraform apply` command in the
terminal.

Once the command completes your new property is created. You can verify
this in [Akamai Control Center](https://control.akamai.com/) or by using the
[Akamai CLI](https://developer.akamai.com/cli). However, you still have
to activate the property configuration.
<!--Is it the Akamai CLI or the PM CLI?-->

## Activate your property

Once the second activation completes, Akamai automatically routes
all traffic to the new property. It deactivates the original
property entirely if all hostnames are no longer pointed at it.

### Create your property activation resource

To activate your property you need to create a new
[akamai_property_activation
resource](https://registry.terraform.io/docs/providers/akamai/r/property_activation).
This resource manages property activations, letting you specify the
property version to activate and the network to activate it on.

You need to set these arguments for this resource:

* property ID and version, which you can set from the `akamai_property` resource.

* network to `STAGING` or `PRODUCTION`. You should activate the property on Staging first to verify that everything works as expected before activating on production. 

* The email addresses to send activation updates to.

* activate argument to `true` to kick off the activation process
<!--I don't see this argument in the example below, or in the akamai_property_activation doc.-->

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


Like you did with the property, you should first verify the
[akamai_property_activation](https://registry.terraform.io/docs/providers/akamai/r/property_activation)
resource by running the `terraform plan` command:

The plan command adds the activation to the property. It doesn't change any other property settings.

If everything looks good, run `terraform apply` to start the activation:

This will activate the property on the selected network.

## How you can use Provisioning resources
<!--Left off editing here.-->

The sections that follow include information on different ways to use
Provisioning resources.

## Dynamic Rule Trees Using Templates

## Dynamic Rule Trees Using Templates

If you wish to inject Terraform interpolations into your rules.json, for example an origin address, you should use rules 
templates to do so.  akamai_property_rules_template data source has many benefits including recursive variable 
replacement and the ability to directly consume Property Manager CLI snippet templates:


More advanced users want different properties to use different rule sets. This can be done by maintaining a base rule 
set and then importing individual rule sets. To do this we first create a directory structure - something like:

```dir
rules/main.json
rules/snippets/routing.json
rules/snippets/performance.json
…
```

The "rules" directory contains a single file "main.json" and a sub directory containing all rule snippets. Here, we 
would provide a basic template for our json.

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

This enables our rules template to process the rules.json and pull each fragment that's referenced. The rendered output 
can be set into property definition resources.

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

## Snippets with Terraform (placeholder)
<!--Is this section still valid? It's from the README.md file.-->

If you wish to inject Terraform interpolations into your rules.json, for example an origin address, you should use rules 
templates to do so.  akamai_property_rules_template data source has many benefits including recursive variable 
replacement and the ability to directly consume Property Manager CLI snippet templates:


More advanced users want different properties to use different rule sets. This can be done by maintaining a base rule 
set and then importing individual rule sets. To do this we first create a directory structure - something like:

```dir
rules/main.json
rules/snippets/routing.json
rules/snippets/performance.json
…
```

The "rules" directory contains a single file "main.json" and a sub directory containing all rule snippets. Here, we 
would provide a basic template for our json.

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

This enables our rules template to process the rules.json and pull each fragment that's referenced. The rendered output 
can be set into property definition resources.

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

