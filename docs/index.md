---
layout: "akamai"
page_title: "Provider: Akamai"
description: |-
  Learn about the Akamai Terraform Provider
---

# Akamai Terraform Provider

Akamai powers and protects life online. Akamai's Intelligent Edge Platform makes it easier for developers and businesses to build, run, and secure applications. We keep decisions, apps, and experiences closer to users, and attacks and threats far away. Our portfolio includes edge security, web and mobile performance, enterprise access, and video delivery solutions.

Use the Akamai Terraform Provider to manage and provision your Akamai configurations in Terraform. You can use the Akamai Provider for many Akamai products. 


## Workflows

Here are the most common workflows for the Akamai Provider:

* **Set up the Provider the first time.** To do this, finish reviewing this guide, then go to [Get Started with the Akamai Provider](guides/get_started_provider.md). When setting up the Provider, you need to choose an [authentication method](guides/akamai_provider_auth.md), and decide whether to import existing Akamai configurations, or create new ones.
* **Add a new module to your existing Akamai Provider configuration.** If the Akamai Provider is already set up, and you're adding a new module, select from the Guides category for the module you're adding.

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
