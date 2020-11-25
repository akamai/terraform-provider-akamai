---
layout: "akamai"
page_title: "Provider: Akamai"
description: |-
  Akamai
---

# Akamai Terraform Provider

Use the Akamai Terraform Provider to manage and provision your Akamai
configurations in Terraform. You can use the Akamai Provider today for 
your Property Manager, Application Security, Edge DNS, and Global
Traffic Management configurations.

Last updated: November 2020.

## Contents

[[Workflows]{.ul}](#workflows)

[[Maintenance]{.ul}](#maintenance)

> [[Manage changes to your Akamai configurations]{.ul}](#manage-changes-to-your-akamai-provider-configurations)
>
> [[Migrate to the newest version of the Akamai Terraform Provider
> (placeholder)]{.ul}](#migrate-to-the-newest-version-of-the-akamai-terraform-provider-placeholder)

[[Links to resources]{.ul}](#links-to-resources)

> [[New to Akamai?]{.ul}](#new-to-akamai)
>
> [[New to Terraform?]{.ul}](#new-to-terraform)

[[Available guides]{.ul}](#available-guides)

## Workflows

Here are the most common workflows for the Akamai Terraform Provider:

* **Set up the Provider the first time.** To do this, finish reviewing this guide, then go to [Set up the Akamai Provider](https://docs.google.com/document/d/1ohurENF2epbu_Dx8X0fcYwAVU1ZEKh\--Wg1jV2AWuP4/edit?usp=sharing). When setting up the Provider, you need to choose an [authentication method](https://docs.google.com/document/d/1S39MM1sZNoM4EmlSLlPVYNohiH6x-Js0IoadUhU4vcc/edit), and decide whether to import existing Akamai configurations, or create new ones.
*  **Add a new module to an existing Provider configuration.** If the Akamai Terraform Provider is already set up, and you're adding a new module, read the guide for the module you're adding:
                              
  * [Application Security](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_appsec)
  * [DNS Zone Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone)
  * [Global Traffic Management](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain)
  * [Property Manager\](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property)
<!-- Need to know the final name for this module. -->  
*  **Update settings for an existing module.** Use the reference information for resource and data sources listed under each module, like DNS or Common. You can find this documentation on the panel to the left.

## Maintenance

### Manage changes to your Akamai configurations

When you're using the Provider, you need to keep your Terraform settings up-to-date 
with any changes made using Akamai APIs, CLIs, and Control Center. You should review 
your network management processes, and update them as needed.

For example, you may want to run `terraform plan` before updating  
your Akamai Provider configurations. You'll likely receive warnings
and suggested changes. Once you any inconsistencies, you can run `terraform plan` 
again and make sure everything is in sync.

### Migrate to the newest version of the Akamai Terraform Provider (placeholder)

<!-- **FOR REVIEWERS: KEEPING THE OLD SECTION BELOW AS A PLACEHOLDER.\
This will likely need to be a separate guide, but we'll need a quick overview here.** -->

<!--**The text below is from "Migrating a property to Terraform" section,
currently in the FAQ doc.** -->

<!--If you have an existing property you would like to migrate to Terraform: -->

<!--1.  Export your rules.json from your existing property using the API,
    > CLI, or Control Center -->

<!--2.  Create a Terraform configuration that pulls in the rules.json -->

<!--3.  Assign a temporary hostname for testing. You can use the edge
    > hostname as the public hostname to allow testing without changing
    > any DNS. -->

<!--4.  Activate the property and test thoroughly. -->

<!--5.  Once testing has concluded successfully, update the configuration to
    > assign the production public hostnames. -->

<!--6.  Activate again.  -->

<!--After this second activation is complete, Akamai automatically routes
all traffic to the new property and deactivates the original property
entirely if no hostnames are pointed at it. -->

## Links to resources

Here are some links to resources that can help get you started with the
Akamai Terraform Provider.

### New to Akamai?

If you're new to Akamai, here are some links to help you get started:

* [Akamai Terraform Provider blogs](https://developer.akamai.com/blog/terraform)
* [Get Started with Akamai APIs](https://developer.akamai.com/api/getting-started)<!--May want a different link.-->
* [Akamai Community site](https://community.akamai.com/customers/s/)

### New to Terraform?

If you're new to Terraform, here are some links to help you get started:

* [A Terraform Tutorial: Download to Installation, to Using Terraform](https://www.terraform.io/downloads.html)
* [A Brief Primer on Terraform's Configuration Language](https://www.terraform.io/docs/configuration/index.html)
* [If you want to learn about Terraform modules](https://www.terraform.io/docs/modules/index.html)
* [Terraform Glossary](https://www.terraform.io/docs/glossary.html)

## Available guides

* [Frequently Asked Questions](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/faq)<!--should be going away-->
* [Get Started with DNS Zone Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone)
* [Get Started with GTM Domain Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain)
* [Get Started with Property Management](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property) <!--Name may change-->
* [Appendix](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/appendix) <!--should be going away-->
