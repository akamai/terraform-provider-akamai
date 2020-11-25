---
layout: "akamai"
page_title: "Akamai: Get Started with the Akamai Terraform Provider"
description: |-
  Learn how to set up the Akamai Terraform Provider for the first time.
---

# Get Started with the Akamai Terraform Provider
<!--Not sure about the name of this doc. -->

If you've set up Akamai APIs before, some of the Akamai Terraform Provider 
setup tasks will look familiar. You'll need to create Akamai API clients 
for each of the modules you'll be using, and retrieve IDs for your contracts 
and groups. Other tasks, like setting up your akamai.tf file, are very
specific to Terraform.

Complete the tasks in this guide when setting up the Akamai
Provider for the first time.

## Contents
<!--Probably don't need TOCs in the final version. Can see TOC in the right pane.-->

* [Get Started](#get-started)

* [Make some decisions](#make-some-decisions)

* [Set up your Terraform folder and configuration file](#set-up-your-terraform-folder-and-configuration-file)

* [Create Akamai API clients](#create-akamai-api-clients)

* [Retrieve contract and group IDs](#retrieve-contract-and-group-ids)

  * [About Groups](#about-groups)

* [Set up your Akamai configurations in Terraform](#set-up-your-akamai-configurations-in-terraform)

* [Initialize the Akamai Provider](#initialize-the-akamai-provider)

* [Test your configuration](#test-your-configuration)

## Get Started 

<!--Not sure this is the right heading for this section.-->
To set up the Akamai Provider, you need to:

* Make some decisions
* Set up your Terraform folder and configuration file
* Create Akamai API clients for each module you'll use.
* Retrieve contract and group IDs
* Set up your Akamai configuration in Terraform
* Initialize the Akamai Provider
* Test your Akamai Provider configuration

## Make some decisions

Before getting into the actual set up of the Akamai Provider, you need
to make some decisions about how you want things to work. You need to
answer these questions:

* **Authentication.** Which type of authentication method do you want to use? Options include:

  * **Local, shared API client.** Uses an account-level API client that all users can access.
  
  * **Individual API client.** Each user on your team needs to set up their own local .edgerc file with their own credentials.
  
  * **Inline credentials.** Have users add their credentials inline when using resources and data sources.
  
  * **Environment variables.** Use environment variables to set credentials. Any variables you set take precedence over the contents of the .edgerc configuration file.

  For details, see *Authenticate your Akamai Terraform Provider*.

* **Modules.** Which modules are you using? The API clients you set up will depend on the modules you choose. For example, if you want to use the Common data sources and resources, you'll need read access to the Property Manager API. If you're not setting up properties in Terraform, you can also obtain contract and group information from Control Center or other Akamai APIs.

* **Akamai configurations.** Are you going to use existing properties and other Akamai configurations with Terraform? Or are you going to start from scratch?

* **Supporting processes.** Are other people in your organization used to making changes via Control Center, an Akamai API, or an Akamai CLI? If they are, you'll need to develop new processes to make sure your Terraform configuration files are fully up to date and the single source of truth.

## Set up your Terraform folder and configuration file
<!--Need to shorten this heading.-->

Now that you made some decisions, you need to: 

1. Create a new folder called `terraform`.
2. Create a file inside your new folder and name it `akamai.tf`.
3. Continue with [Create Akamai API clients](#create-akamai-api-clients), where you'll create Akamai API clients and add credential information to your `akamai.tf` file.

## Create Akamai API clients

Create an Akamai API Client with the right permissions and valid
credentials to authenticate your Akamai Terraform files. Your Akamai API
Client needs read-write permission to the APIs associated with the
Akamai Provider modules you're using, like DNS or Provisioning.

Refer to the [Authenticate Your Akamai Provider](https://docs.google.com/document/d/1S39MM1sZNoM4EmlSLlPVYNohiH6x-Js0IoadUhU4vcc/edit\#heading=h.f8au8xqqw0yw)
guide for details.
<!--Need final link.-->

Once you're done authenticating, come back here to complete the rest of
the Akamai Provider setup.

## Retrieve contract and group IDs

You'll need contract and group IDs to use most of the Akamai Terraform
Provider modules. These IDs provide needed account information.

You can retrieve these IDs through the `akamai_contract` and
`akamai_group` data sources, which require read access to the Property
Manager API. You can also get this information from the Contracts app in
Control Center, or by using other Akamai APIs.

### About Groups

If you're not familiar with Akamai groups, they control access to your
Akamai configurations and help consolidate reporting functions. Each account
 features a hierarchy of groups, which typically map to an organizational hierarchy.

Your account admins can use Control Center or the [Identity Management: User Administration API](https://developer.akamai.com/en-us/api/core_features/identity_management_user_admin/v2.html)
to set up groups, each with their own set of users and roles.

## Set up your Akamai configurations in Terraform

You're now ready to import existing configurations or create new ones
from scratch.

At this point in the setup, you should refer to the guides for the
Akamai modules you're using:

| **Module** | **Guide** |
|------------|------------|
| Application Security | [Get Started with Application Security](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_appsec) |
| Edge DNS (DNS) | [Get Started with DNS Zone Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone) | 
| Global Traffic Management | [Get Started with GTM Domain Administration](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain) | 
| Property Manager (Provisioning and Common modules) | [Get Started with Property Management](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property) |

-> **Note** Both Terraform and the Akamai Terraform CLI package come
pre-installed in the Akamai Development Environment. Get more details in
our [[installation
Instructions](https://developer.akamai.com/blog/2020/05/26/set-development-environment).

Once you're done with the module-level setup, continue with the next
sections to initialize and test the Akamai Provider.

## Initialize the Akamai Provider
<!--May want to put this sections in the individual guides.-->

Once you have your configuration complete, save the `akamai.tf` file. Then
switch to the terminal to initialize Terraform using the command:

`$ terraform init`

This command installs the latest version of the Akamai Provider, as well
as any other providers you're using. To update
the Akamai provider version after a new release, simply run terraform
init again.

## Test your configuration
<!--May want to put this sections in the individual guides.-->

To test your configuration, use the plan command:

`$ terraform plan`

This command makes Terraform create a plan for the work it will do
based on the `akamai.tf` configuration file. It doesn't actually make any changes
and is safe to run as many times as you like.