---
layout: "akamai"
page_title: "Provider: Akamai"
description: |-
  Learn about the Akamai Terraform Provider
---

# Akamai Terraform Provider

Use the Akamai Terraform Provider to manage and provision your Akamai configurations in Terraform. 
You can use the Akamai Provider today for many Akamai products.

!> Version 1.0.0 of the Akamai Terraform Provider is a major release that's currently available for the Provisioning module. Before upgrading, you need to make changes to some of your Provisioning resources and data sources. See [Upgrade to Version 1.0.0](guides/1.0_migration.md) for details.

Last updated: February 2022.

## Migrate to the newest version

If you're using the Provisioning module, the latest major version of the Akamai Provider is now available. 
See [Upgrade to Version 1.0.0](guides/1.0_migration.md) for more information.

## Workflows

Here are the most common workflows for the Akamai Provider:

* **Set up the Provider the first time.** To do this, finish reviewing this guide, then go to [Get Started with the Akamai Provider](guides/get_started_provider.md). When setting up the Provider, you need to choose an [authentication method](guides/akamai_provider_auth.md), and decide whether to import existing Akamai configurations, or create new ones.
* **Add a new module to your existing Akamai Provider configuration.** If the Akamai  Provider is already set up, and you're adding a new module, read the guide for the module you're adding:
                              
  * [Application Security](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_appsec)
  * [Certificate Provisioning (CPS)](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_cps)
  * [Cloudlets](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_cloudlets)
  * [DataStream](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_datastream)
  * [DNS Zone Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone)
  * [EdgeWorkers](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_edgeworkers)
  * [Global Traffic Management Domain Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain)
  * [Identity and Access Management](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_iam.md)
  * [Network Lists](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_networklists)
  * [Provisioning/Property Manager](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property)
*  **Update settings for an existing module.** Use the reference information for resource and data sources listed under each module, like DNS Zone Administration or Property Provisioning. You can find this documentation on the panel to the left.

## Manage changes to your Akamai configurations

When you're using the Akamai Provider, you need to keep your Terraform configurations up to date with changes made using Akamai APIs, CLIs, and Control Center. 
You should review your network management processes and update them to include the Akamai Provider.

For example, before updating your Akamai Provider configurations, you may want to run `terraform plan` first. 
You'll likely receive warnings and suggested changes. 
Once you fix any issues, you can run `terraform plan` again and make sure everything is in sync.


## Links to resources

Here are some links to resources that can help get you started with the Akamai Terraform Provider.

### New to Akamai?

If you're new to Akamai, here are some links to help you get started:

* [Create authentication credentials](https://techdocs.akamai.com/developer/docs/set-up-authentication-credentials)
* [Akamai Community site](https://community.akamai.com/customers/s/)

### New to Terraform?

If you're new to Terraform, here are some links you might find helpful:

* [Get started tutorials](https://learn.hashicorp.com/collections/terraform/cloud-get-started)
* [Terraform module tutorials](https://learn.hashicorp.com/collections/terraform/modules)
* [Terraform configuration language tutorials](https://learn.hashicorp.com/collections/terraform/configuration-language)
* [Terraform glossary](https://www.terraform.io/docs/glossary.html)

## Available guides

Here's a list of the guides for the Akamai Provider in the general order you might use them:

* [Get Started with the Akamai Terraform Provider](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_provider)
* [Authenticate the Akamai Terraform Provider](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/akamai_provider_auth)
* [Upgrade to Version 1.0.0](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/1.0_migration)
* [Application Security Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_appsec)
* [Certificate Provisioning Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_cps)
* [Cloudlets Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_cloudlets)
* [DataStream Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_datastream)
* [EdgeWorkers Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_edgeworkers)
* [DNS Zone Administration Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone)
* [Global Traffic Management Domain Administration Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain)
* [Identity and Access Management Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_iam.md)
* [Network Lists Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_networklists)
* [Property Provisioning Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property)
* [Shared Resources](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/appendix) 
