---
layout: "akamai"
page_title: "Akamai: FAQ (Frequently Asked Questions)"
sidebar_current: "docs-akamai-guide-faq"
description: |-
  Frequently Asked Questions
---

# Frequently Asked Questions

## Migrating a property to Terraform

If you have an existing property you would like to migrate to Terraform we recommend the following process:

1. Export your rules.json from your existing property (using the API, CLI, or Control Center)
2. Create a terraform configuration that pulls in the rules.json
3. Assign a temporary hostname for testing (hint: you can use the edge hostname as the public hostname to allow testing without changing any DNS)
4. Activate the property and test thoroughly
5. Once testing has concluded successfully, update the configuration to assign the production public hostnames
6. Activate again

Once this second activation completes Akamai will automatically route all traffic to the new property and will deactivate the original property entirely if all hostnames are no longer pointed at it.

Since Terraform assumes it is the de-facto state for any resource it leverages, we strongly recommend creating a new property based off an existing rules.json tree when starting with the provider to mitigate any risks to existing setups. 

## Dynamic Rule Trees Using Templates

If you wish to inject terraform data into your rules.json, for example an origin address, you can use Terraform templates to do so like so:

First decide where your origin value will come from, this could be another Terraform resource such as your origin cloud provider, or it could be a terraform input variable like this:

```hcl
variable "origin" { }
```

Because we have not specified a default, a value is required when executing the config. We can then reference this variable using `${vars.origin}` in our template data source:

```hcl
data "template_file" "init" {
  template = "${file("rules.json")}"
  vars = {
    origin = "${vars.origin}"
  }
}
```

Then in our `rules.json` we would have:

```json
{
  "name": "origin",
  "options": {
    "hostname": "**${origin}**",
    ...
  }
},
```

You can also inject entire JSON blocks using the same mechanism:

```json
{
	"rules": {
		"behaviors": [
    		${origin}

	    ]
	}
}
```
## Migrating a GTM domain (and contained objects) to Terraform

Migrating an existing GTM domain can be done in many ways. Two such methods include:

### Via Command Line Utility

A package, [CLI-Terraform-GTM](https://github.com/akamai/cli-terraform-gtm), for the [Akamai CLI](https://developer.akamai.com/cli) provides a time saving means to collect information about, generate a configuration for, and import an existing GTM domain and its contained objects and attributes. With the package, you can:

1. Generate a json formatted list of all domain objects
2. Generate a Terraform configuration for the domain and contained objects
3. Generate a command line script to import all defined resources

It is recommended that the existing domain configuration (using the API or Control Center) be downloaded before hand as a backup and reference.  Additionally, a terraform plan should be executed after importing to validate the generated tfstate. Note: The first time plan is run, an update will be shown for the provider defined domain fields: contract, group and wait_on_complete.

### Via Step By Step Construction

1. Download your existing domain configuration (using the API or Control Center) as a backup and reference.
2. Using the domain download as a reference, create a terraform configuration representing the the existing domain and all contained GTM objects. Note: In creating each resource block, make note of `required`, `optional` and `computed` fields.
3. Use the Terraform Import command to import the existing domain and contained objects; singularly and in serial order.
4. (Optional, Recommended) Review domain download content and created terraform.tfstate to confirm the domain and all objects are represented correctly
5. Execute a `Terraform Plan` on the configuration. The plan should be empty. If not, correct accordingly and repeat until plan is empty and configuration is in sync with the GTM Backend.

Since Terraform assumes it is the de-facto state for any resource it leverages, we strongly recommend staging the domain and objects imports in a test environment to familiarize yourself with the provider operation and mitigate any risks to the existing GTM domain configuration.

## Upgrading the Akamai Provider

To upgrade the provider, simply run `terraform init` again, and all providers will be updated to their latest version within specified version constraints.
