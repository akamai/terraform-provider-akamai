---
layout: "akamai"
page_title: "Akamai: Get Started with GTM Domain Administration"
description: |-
  Get Started with Akamai GTM Domain Administration using Terraform
---

# Get Started with GTM Domain Administration

The Akamai Provider for Terraform provides you the ability to automate the creation, deployment, and management of GTM domain configuration and administration; as well as importing existing domains and contained objects.  

To get more information about Global Traffic Management, see:

* Developer - [API documentation](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html)
* User Guide - [Official Documentation](https://learn.akamai.com/en-us/products/web_performance/global_traffic_management.html)

Remember to start by reviewing the Get Started with the [Akamai Terraform Provider Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_provider). You should have a valid `akamai.tf` Terraform configuration at this point before adding the GTM module configuration.

## Create a GTM Domain

The domain itself is represented by a [`akamai_gtm_domain` resource](../resources/gtm_domain.md). Add this new resource block to your `akamai.tf` file after the provider block. Note: the domain must be the first GTM resource created as it provides operating context for all other contained objects.

To define the entire configuration, we start by opening the resource block and giving the domain a name. In this case we’re going to use the name "example".

Next, we set the required (domain, type) and optional (group ID, contract ID, email list, comment) arguments.

Once you’re done, your domain configuration should look like this:

```hcl
resource "akamai_gtm_domain" "example" {
	name = "example.akadns.net"                     # Domain Name
	type = "weighted"				# Domain type
	group_id    = data.akamai_group.default.id         # Group ID variable
	contract_id = data.akamai_contract.default.id      # Contract ID variable
	email_notification_list = ["user@demo.me"]        # email notification list
	comment = "example domain demo"
}
```
> **Note:** Notice that we’re using variables from the previous section to reference the group and contract IDs. These will automatically be replaced at runtime by Terraform with the actual values.

## Create a GTM Datacenter

The datacenter itself is represented by an [`akamai_gtm_datacenter` resource](../resources/gtm_datacenter.md). Add this new block to your `akamai.tf` file after the provider block.

To define the entire configuration, we start by opening the resource block and give it a name. In this case we’re going to use the name "example_dc".

Next, we set the required (domain name) and optional (nickname) arguments.

Once you’re done, your datacenter configuration should look like this:

```hcl
resource "akamai_gtm_datacenter" "example_dc" {
	domain = akamai_gtm_domain.example.name		# domain
	nickname = "datacenter_1"   			# Datacenter Nickname
	depends_on = [akamai_gtm_domain.example]
}
```

## Create a GTM Property

The property itself is represented by an [`akamai_gtm_property` resource](../resources/gtm_property.md). Add this new block to your `akamai.tf` file after the provider block.

To define the entire configuration, we start by opening the resource block and give it a name. In this case we’re going to use the name "example_prop".

Next, we set the required (domain name, property name, property type, traffic_targets, liveness_tests, score_aggregation_type, handout_limit, handout_mode) and optional (failover_delay, failback_delay) arguments.

Once you’re done, your property configuration should look like this:

```hcl
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

This command will make Terraform create a plan for the work it will do based on the configuration file. This will not actually make any changes and is safe to run as many times as you like.

## Apply Changes

To actually create our domain, datacenter and property;, we need to instruct Terraform to apply the changes outlined in the plan. To do this, in the terminal, run the command:

```
$ terraform apply
```

Once this completes your domain, datacenter and property will have been created. You can verify this in [Akamai Control Center](https://control.akamai.com) or via the [Akamai CLI](https://developer.akamai.com/cli).

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

[Migrating A GTM Domain](../guides/faq.md#migrating-a-gtm-domain-and-contained-objects-to-terraform) discusses GTM resource import in more detail.
