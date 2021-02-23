---
layout: "akamai"
page_title: "Provider: Akamai"
description: |-
  Learn about the Akamai Terraform Provider
---

# Akamai Terraform Provider

Use the Akamai Terraform Provider to manage and provision your Akamai
configurations in Terraform. You can use the Akamai Provider today for 
your Property Manager, Application Security, Edge DNS, and Global
Traffic Management configurations.

!> Version 1.0.0 of the Akamai Terraform Provider is a major release that's currently available for the Provisioning module. Before upgrading, you need to make changes to some of your Provisioning resources and data sources. See the [migration guide](guides/1.0_migration.md) for details.

Last updated: February 2021.

## Migrate to the newest version

If you're using the Provisioning module, the latest major version of the Akamai Provider is now available. See [1.0.0 Migration Guide](guides/1.0_migration.md) for more information.

## Workflows

Here are the most common workflows for the Akamai Provider:

* **Set up the Provider the first time.** To do this, finish reviewing this guide, then go to [Get Started with the Akamai Terraform Provider](guides/get_started_provider.md). When setting up the Provider, you need to choose an [authentication method](guides/akamai_provider_auth.md), and decide whether to import existing Akamai configurations, or create new ones.
* **Add a new module to your existing Akamai Provider configuration.** If the Akamai  Provider is already set up, and you're adding a new module, read the guide for the module you're adding:
                              
  * [Application Security](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_appsec)
  * [DNS Zone Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone)
  * [Global Traffic Management](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain)
  * [Identity and Access Management](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_iam.md)
  * [Provisioning/Property Manager](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property)
*  **Update settings for an existing module.** Use the reference information for resource and data sources listed under each module, like DNS or Provisioning. You can find this documentation on the panel to the left.

## Manage changes to your Akamai configurations

When you're using the Akamai Provider, you need to keep your 
Terraform configurations up to date with changes made using Akamai 
APIs, CLIs, and Control Center. You should review your network management 
processes and update them to include the Akamai Provider.

For example, before updating your Akamai Provider configurations, you may want to
 run `terraform plan` first. You'll likely receive warnings
and suggested changes. Once you fix any issues, you can run `terraform plan` 
again and make sure everything is in sync.


## Links to resources

Here are some links to resources that can help get you started with the
Akamai Terraform Provider.

### New to Akamai?

If you're new to Akamai, here are some links to help you get started:

* [Get Started with Akamai APIs](https://developer.akamai.com/api/getting-started)
* [Akamai Community site](https://community.akamai.com/customers/s/)

### New to Terraform?

If you're new to Terraform, here are some links you might find helpful:

* [A Terraform Tutorial: Download to Installation, to Using Terraform](https://www.terraform.io/downloads.html)
* [A Brief Primer on Terraform's Configuration Language](https://www.terraform.io/docs/configuration/index.html)
* [If you want to learn about Terraform modules](https://www.terraform.io/docs/modules/index.html)
* [Terraform Glossary](https://www.terraform.io/docs/glossary.html)

## Available guides

Here's a list of the guides for the Akamai Provider in the general order you might use them:

* [Get Started with the Akamai Terraform Provider](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_provider)
* [Authenticate the Akamai Terraform Provider](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/akamai_provider_auth)
* [Akamai Terraform Provider: 1.0.0 Migration Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/1.0_migration)
* [Get Started with the Identity and Access Management Module](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_iam.md)
* [Get Started with the Provisioning Module](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property)
* [Get Started with Application Security](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_appsec)
* [Get Started with DNS Zone Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone)
* [Get Started with GTM Domain Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain)
*  [Appendix](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/appendix) 
