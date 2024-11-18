---
layout: "akamai"
page_title: "Provider: Akamai"
description: |-
  Learn about the Akamai Terraform Provider
---

~> We’ve moved our subprovider documentation. You can now find it in the [Terraform guide](https://techdocs.akamai.com/terraform/docs/overview) on Akamai TechDocs.

# Akamai Terraform Provider

Akamai powers and protects life online. Akamai's Intelligent Edge Platform makes it easier for developers and businesses to build, run, and secure applications. We keep decisions, apps, and experiences closer to users, and attacks and threats far away. Our portfolio includes edge security, web and mobile performance, enterprise access, and video delivery solutions.

Use the Akamai Terraform Provider to manage and provision your Akamai configurations in Terraform. You can use the Akamai Provider for many Akamai products. 

## Workflows

Here are the most common workflows for the Akamai Provider:

* **Set up the Provider the first time.** To do this, finish reviewing this guide, then go to [Get Started with the Akamai Provider](https://techdocs.akamai.com/terraform/docs/overview). When setting up the Provider, you need to choose an [authentication method](https://techdocs.akamai.com/terraform/docs/overview#add-authentication), and decide whether to import existing Akamai configurations, or create new ones.
* **Add a new module to your existing Akamai Provider configuration.** If the Akamai Provider is already set up, and you're adding a new module, select from the Guides category for the module you're adding.

## Akamai configuration management

When you're using the Akamai Provider, you need to keep your Terraform configurations up to date with changes made using Akamai APIs, CLIs, and Control Center. 
You should review your network management processes and update them to include the Akamai Provider.

For example, before updating your Akamai Provider configurations, you may want to run `terraform plan` first. 
You'll likely receive warnings and suggested changes. 
Once you fix any issues, you can run `terraform plan` again and make sure everything is in sync.

## Subprovider documentation

We’ve moved our documentation to the Akamai TechDocs site. Use the table to find information about the subprovider you’re using.

| Subprovider                                                                                  | Description                                                                                          |
|----------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------|
| [Application Security](https://techdocs.akamai.com/terraform/v6.6/docs/configure-appsec)     | Manage security configurations, security policies, match targets, rate policies, and firewall rules. |
| [Bot Manager](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-botman)                 | Identify, track, and respond to bot activity on your domain or in your app.                          |
| [Certificates](https://techdocs.akamai.com/terraform/v6.6/docs/cps-integration-guide)        | Full life cycle management of SSL certificates for your ​Akamai​ CDN applications.                   |
| [Client Lists](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-client-lists)          | Reduce harmful security attacks by allowing only trusted IP/CIDRs, locations, autonomous system numbers, and TLS fingerprints to access your services and content.|
| [Cloud Access Manager](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-cam)           | Enable cloud origin authentication and securely store and manage your cloud origin credentials as access keys. |
| [Cloud Wrapper](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-cloud-wrapper)        | Provide your customers with a more consistent user experience by adding a custom caching layer that improves the connection between your cloud infrastructure and the Akamai platform.|
| [Cloudlets](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-cloudlets)                | Solve specific business challenges using value-added apps that complement ​Akamai​'s core solutions. |
| [DataStream](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-datastream)              | Monitor activity on the ​Akamai​ platform and send live log data to a destination of your choice.    |
| [Edge DNS](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-edgedns)                   | Replace or augment your DNS infrastructure with a cloud-based authoritative DNS solution.            |
| [EdgeWorkers](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-edgeworkers)            | Execute JavaScript functions at the edge to optimize site performance and customize web experiences. |
| [Global Traffic Management](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-gtm)      | Use load balancing to manage website and mobile performance demands.                                 |
| [Identity and Access Management](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-iam) | Create users and groups, and define policies that manage access to your Akamai applications.         |
| [Image and Video Manager](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-ivm)        | Automate image and video delivery optimizations for your website visitors.                           |
| [Network Lists](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-network-lists)        | Automate the creation, deployment, and management of lists used in ​Akamai​ security products.       |
| [Property](https://techdocs.akamai.com/terraform/v6.6/docs/set-up-property-provisioning)     | Define rules and behaviors that govern your website delivery based on match criteria.                |

## Links to resources

Here are some links to resources to help you get started:

* [Create Akamai authentication credentials](https://techdocs.akamai.com/terraform/docs/overview#add-authentication)
* [Akamai's Terraform Tapas video series](https://www.youtube.com/playlist?list=PLDlttLRccCk7a-JNb-xFH6dz4WqG53JQa)
* [Akamai Community site](https://community.akamai.com/customers/s/)
* [Terraform tutorials](https://learn.hashicorp.com/collections/terraform/cloud-get-started)
* [Terraform module tutorials](https://learn.hashicorp.com/collections/terraform/modules)
* [Terraform configuration language tutorials](https://learn.hashicorp.com/collections/terraform/configuration-language)
* [Terraform glossary](https://www.terraform.io/docs/glossary.html)