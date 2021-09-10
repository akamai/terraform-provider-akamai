---
layout: "akamai"
page_title: "Module: GTM Domain Administration"
description: |-
  GTM Domain Administration module for the Akamai Terraform Provider
---

# Global Traffic Management Domain Administration Module Guide

The Akamai Provider for Terraform provides you the ability to automate the creation, deployment, and management of Global Traffic Management (GTM) domain configuration and administration; as well as importing existing domains and contained objects.  

To get more information about Global Traffic Management (GTM), see:

* Developer - [API documentation](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html)
* User Guide - [Official Documentation](https://learn.akamai.com/en-us/products/web_performance/global_traffic_management.html)

## Prerequisites 

Before starting with the DNS module, you need to:

1. Complete the tasks in [Get Started with the Akamai Provider](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_provider). You should have an API client and a valid `akamai.tf` Terraform configuration before adding the GTM module configuration.
2. Determine whether you want to import an existing DNS zone and records or create new ones.
3. If you're importing an existing GTM domain, continue with [Import a GTM  domain](#import-a-gtm-domain).
3. If you're creating a new GTM domain, continue with [Create a GTM domain
](#create-a-gtm-domain).

## Import a GTM domain

You can migrate an existing GTM domain into your Terraform configuration using either a command line utility or step-by-step construction.

### Import using the command line utility

You can use the [Akamai CLI for Akamai Terraform Provider](https://github.com/akamai/cli-terraform) to generate a configuration for and import an existing GTM domain. With the package, you can generate:

* a JSON-formatted list of all domain objects.
* a Terraform configuration for the domain and contained objects.
* a command line script to import all defined resources.

Before using this CLI, keep the following in mind:

* Download the existing GTM domain configuration to have as a backup and reference during an import. You can download these by using the [Global Traffic Management API](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html) or the GTM app on [Control Center](https://control.akamai.com).  
* Terraform limits the characters that can be part of its resource names. During construction of the resource configurations, invalid characters are replaced with underscore , '_'.
* Terraform doesn't provide any state information during import. When you run `plan` and `apply` after an import, Terraform lists discrepencies and reconciles configurations and state. Any discrepencies clear following the first `apply`. 
* After first time you run `plan` or `apply`, the `contract`, `group`, and `wait_on_complete` attributes are updated.
* Run `terraform plan` after importing to validate the generated `tfstate` file.

### Import using step-by-step construction

To import using step-by-step construction, complete these tasks:

1. Determine how you want to test your Terraform import. For example, you may want to set up your zone and recordset imports in a test environment to familiarize yourself with the provider operation and mitigate any risks to your existing DNS zone configuration.
1. Download the existing domain configuration and master file to have as a backup and reference during an import. You can download these from the [Global Traffic Management API](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html) or from the GTM app on [Control Center](https://control.akamai.com) .  
1. Using the domain download as a reference, create a Terraform configuration representing the existing domain and all contained GTM objects. 
1. Verify that your Terraform configuration addresses all required attributes and any optional and computed attributes you need.
1. Run `terraform import`. This command imports the existing domain and contained objects one at a time based on the order in the configuration.
1. Compare the downloaded domain file with the `terraform.tfstate` file to confirm that the domain and all objects are represented correctly.
1. Run `terraform plan` on the configuration. The plan should be empty. If not, correct accordingly and repeat until plan is empty and configuration is in sync with the Edge DNS backend.### Via Step By Step Construction
1. Run `terraform plan` on the configuration. The plan should be empty. If not, correct accordingly and repeat until plan is empty and configuration is in sync with the GTM Backend.

## Create a GTM Domain

The Domain itself is represented by a [`akamai_gtm_domain` resource](../resources/gtm_domain.md). Add this new resource block to your `akamai.tf` file after the provider block. **Note:** the domain must be the first GTM resource created as it provides operating context for all other contained objects.

To define the entire configuration, we start by opening the resource block and giving the domain a `name`. In this case, we're going to use the name "example".

Next, we set the required (`name`, `type`) and optional (`group_id`, `contract_id`, `email_notification_list`, `comment`) arguments.

Once you're done, your Domain configuration should look like this:

```
resource "akamai_gtm_domain" "example" {
	name = "example.akadns.net"                     # Domain Name
	type = "weighted"				# Domain type
	group_id    = data.akamai_group.default.id         # Group ID variable
	contract_id = data.akamai_contract.default.id      # Contract ID variable
	email_notification_list = ["user@demo.me"]        # email notification list
	comment = "example domain demo"
}
```
> **Note:** Notice the use of variables from the previous section to reference the group and contract IDs. These will be replaced at runtime by Terraform with the actual values.

## Create a GTM Datacenter

The Datacenter itself is represented by a [`akamai_gtm_datacenter` resource](../resources/gtm_datacenter.md). Add this new block to your `akamai.tf` file after the provider block.

To define the entire configuration, we start by opening the resource block and giving it a name. In this case, we're going to use the name "example_dc".

Next, we set the required (`domain` name) and optional (`nickname`) arguments.

Once done, your Datacenter configuration should look like this:

```
resource "akamai_gtm_datacenter" "example_dc" {
	domain = akamai_gtm_domain.example.name		# domain
	nickname = "datacenter_1"   			# Datacenter Nickname
	depends_on = [akamai_gtm_domain.example]
}
```

## Create a GTM Property

The Property itself is represented by a [`akamai_gtm_property` resource](../resources/gtm_property.md). Add this new block to your `akamai.tf` file after the provider block.

To define the entire configuration, we start by opening the resource block and giving it a name. In this case, we're going to use the name "example_prop".

Next, we set the required (`domain` name, property `name`, property `type`, `traffic_target`s, `liveness_test`s, `score_aggregation_type`, `handout_limit`, `handout_mode`) and optional (`failover_delay`, `failback_delay`) arguments.

Once you're done, your Property configuration should look like this:

```
resource "akamai_gtm_property" "example_prop" {
	domain = akamai_gtm_domain.example.name         # domain
	name = "example_prop_1"                         # Property Name
	type = "weighted-round-robin"
	score_aggregation_type = "median"
	handout_limit = 5
	handout_mode = "normal"
	failover_delay = 0 
	failback_delay = 0
	traffic_target = {
		datacenter_id = akamai_gtm_datacenter.example_dc.datacenter_id
		enabled = true
		weight = 100
		servers = ["1.2.3.4"]
		name = ""
		handout_cname = ""
	}
	liveness_test = {
		name = "lt1"
		test_interval = 10
		test_object_protocol = "HTTP"
		test_timeout = 20
		answer_required = false
		disable_nonstandard_port_warning = false
		error_penalty = 0
		host_header = ""
		http_error3xx = false
		http_error4xx = false
		http_error5xx = false
		disabled = false
		peer_certificate_verification = false
		recursion_requested = false
		request_string = ""
		resource_type = ""
		response_string = ""
		ssl_client_certificate = ""
		ssl_client_private_key = ""
		test_object = "/junk"
		test_object_password = ""
		test_object_port = 1
		test_object_username = ""
		timeout_penalty = 0
	}
	depends_on = [
		akamai_gtm_domain.example,
		akamai_gtm_datacenter.example_dc
	]
}
```

## Initialize the Provider

Once you have your configuration complete, save the file. Then switch to the terminal to initialize Terraform using the command:

```
$ terraform init
```

This command will install the latest version of the Akamai Provider, as well as any other providers necessary (such as the local provider). To update the Akamai Provider version after a new release, simply run `terraform init` again.

## Test Your Configuration

To test your configuration, use `terraform plan`:

```
$ terraform plan
```

This command will make Terraform create a plan for the work set by the configuration file. This will not actually make any changes and is safe to run as many times.

## Apply Changes

To actually create our Domain, Datacenter and Property, we need to instruct Terraform to apply the changes outlined in the plan. To do this, run the command:

```
$ terraform apply
```

Once this completes your Domain, Datacenter and Property will have been created. You can verify this in [Akamai Control Center](https://control.akamai.com) or via the [Akamai CLI](https://developer.akamai.com/cli).

## Import Existing GTM Resource

Existing GTM resources may be imported using the following formats:

```
$ terraform import akamai_gtm_domain.{{domain resource name}} {{gtm domain name}}
$ terraform import akamai_gtm_datacenter.{{datacenter resource name}} {{gtm domain name}}:{{gtm datacener id}}
$ terraform import akamai_gtm_property.{{property resource name}} {{gtm domain name}}:{{gtm property name}}
$ terraform import akamai_gtm_resource.{{resource resource name}} {{gtm domain name}}:{{gtm resource name}}
$ terraform import akamai_gtm_cidrmap.{{cidrmap resource name}} {{gtm domain name}}:{{gtm cidrmap name}}
$ terraform import akamai_gtm_geomap.{{geomap resource name}} {{gtm domain name}}:{{gtm geographicmap name}}
$ terraform import akamai_gtm_asmap.{{asmap resource name}} {{gtm domain name}}:{{gtm asmap name}}
```

## GTM field status when running plan and apply

When using `terraform plan` or `terraform apply`, Terraform presents fields defined in the configuration and all defined resource fields. Fields are either required, optional, or computed as specified in each resource description. Default values for fields will display if not explicitly configured. In many cases, the default will be zero, empty string, or empty list depending on the type. These default or empty values are informational and not included in resource updates.
