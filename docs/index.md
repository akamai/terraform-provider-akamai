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

Last updated: December 2020.

## Workflows

Here are the most common workflows for the Akamai Provider:

* **Set up the Provider the first time.** To do this, finish reviewing this guide, then go to [Get Started with the Akamai Terraform Provider
](https://docs.google.com/document/d/1ohurENF2epbu_Dx8X0fcYwAVU1ZEKh\--Wg1jV2AWuP4/edit?usp=sharing). When setting up the Provider, you need to choose an [authentication method](https://docs.google.com/document/d/1S39MM1sZNoM4EmlSLlPVYNohiH6x-Js0IoadUhU4vcc/edit), and decide whether to import existing Akamai configurations, or create new ones.
* **Add a new module to your existing Akamai Provider configuration.** If the Akamai  Provider is already set up, and you're adding a new module, read the guide for the module you're adding:
                              
  * [Application Security](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_appsec)
  * [DNS Zone Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone)
  * [Global Traffic Management](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain)
  * [Provisioning/Property Manager](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property)
*  **Update settings for an existing module.** Use the reference information for resource and data sources listed under each module, like DNS or Common. You can find this documentation on the panel to the left.

## Maintenance
<!--Might remove this heading as the two subheads won't appear in the right-hand TOC.-->

### Manage changes to your Akamai configurations

When you're using the Akamai Provider, you need to keep your 
Terraform configurations up-to-date with changes made using Akamai 
APIs, CLIs, and Control Center. You should review your network management 
processes and update them to include the Akamai Provider.

For example, before updating your Akamai Provider configurations, you may want to
 run `terraform plan` first. You'll likely receive warnings
and suggested changes. Once you fix any issues, you can run `terraform plan` 
again and make sure everything is in sync.

### Migrate to the newest version of the Akamai Provider 

<!-- This section is a placeholder for the migration guide being developed. Likely need some overview text and a link to a separate migration guide.-->

## Links to resources

Here are some links to resources that can help get you started with the
Akamai Terraform Provider.

### New to Akamai?

If you're new to Akamai, here are some links to help you get started:

* [Akamai Provider blogs](https://developer.akamai.com/blog/terraform)
* [Get Started with Akamai APIs](https://developer.akamai.com/api/getting-started)<!--May want a different link.-->
* [Akamai Community site](https://community.akamai.com/customers/s/)

### New to Terraform?

If you're new to Terraform, here are some links you might find helpful:

* [A Terraform Tutorial: Download to Installation, to Using Terraform](https://www.terraform.io/downloads.html)
* [A Brief Primer on Terraform's Configuration Language](https://www.terraform.io/docs/configuration/index.html)
* [If you want to learn about Terraform modules](https://www.terraform.io/docs/modules/index.html)
* [Terraform Glossary](https://www.terraform.io/docs/glossary.html)

## Available guides

* [Frequently Asked Questions](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/faq)
<!--May be going away with Dec release. More likely will stay until next doc update.-->
* [Get Started with DNS Zone Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone)
* [Get Started with GTM Domain Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain)
* [Get Started with Property Management](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property) <!--Name may change-->
* [Appendix](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/appendix) 
<!--Want to rename to something like "Common codes and formats". Any suggestions?-->
