---
layout: "akamai"
page_title: "Akamai Provider: Get Started"
description: |-
  Learn how to set up the Akamai Terraform Provider for the first time.
---

# Get Started with the Akamai Provider

Complete the tasks in this guide when setting up the Akamai
Provider for the first time.

If you've set up Akamai APIs before, some of the Akamai Provider
setup tasks will look familiar. You need to create Akamai API clients
for each of the modules you'll be using, and retrieve IDs for your contracts
and groups. Other tasks, like setting up your `.tf` configuration file, are
specific to Terraform.

-> **Note** Both Terraform and the Akamai Terraform CLI package come
pre-installed in the Akamai Development Environment. See [Set Up a Development Environment](https://developer.akamai.com/blog/2020/05/26/set-development-environment) for more information.

## Workflow

To set up the Akamai Provider, you need to:

* Make initial decisions about how you want to work with the Akamai Provider.
* Set up your Terraform folder and configuration file.
* Create Akamai API clients for each module you'll use.
* Retrieve contract and group IDs.
* Arrange resources and data sources in the Akamai configuration file.
* Initialize the Akamai Provider.
* Test your Akamai Provider configuration.

## Make initial decisions

Before getting into the actual set up, you need
to decide how you want to work with the Akamai Provider. You need to
answer these questions:

* **Authentication.** Which type of authentication method do you want to use? Options include:

  * **Local, shared API client.** Uses an account-level API client that all users can access.

  * **Individual API client.** Each user on your team needs to set up their own local `.edgerc` file with their own credentials.

  * **Inline credentials.** Have users add their credentials inline when using resources and data sources.

  * **Environment variables.** Use environment variables to set credentials. Any variables you set take precedence over the contents of the `.edgerc` configuration file.

  For details, see [Authenticate the Akamai Provider](../guides/akamai_provider_auth.md).

* **Modules.** Which modules do you plan to use? The API clients you set up depend on the modules you choose. For example, if you want to use the Common data sources and resources, you'll need read access to the Property Manager API.

* **Akamai configurations.** Are you going to use existing properties and other Akamai configurations with Terraform? Or are you going to start from scratch?

* **Single or multiple `.tf` files.** Do you want to manage the full lifecycle of your infrastructure in a single file? Or do you prefer to split it into smaller Terraform configurations with limited scope and delegate them to specific teams? Independent configurations use output variables to publish information and enable access to that data from other workspaces.

* **Supporting processes.** Are other people in your organization used to making changes via Control Center, an Akamai API, or an Akamai CLI? If they are, you need to develop new processes to make sure your Terraform configuration files are fully up to date and from now on, the single source of changes. All the modifications your team makes outside of Terraform get overwritten whenever you [run the `terraform -apply` command](#apply-your-configuration).

## Set up your `.tf` files

Now that you have all the answers, set up a Terraform configuration files for the Akamai modules you plan to use.

1. Create a new folder called `terraform`.
2. Create a file inside your new folder and name it `akamai.tf`. If you decided to split the configuration into smaller chunks, create all the files accordingly.
3. Continue with [Create Akamai API clients](#create-akamai-api-clients).

## Create Akamai API clients

Create an Akamai API client with the right permissions and valid credentials to authenticate your Akamai Provider files. The Akamai API client needs read-write permission to the APIs associated with the Akamai Provider modules you're using, like DNS Zone Administration or Property Provisioning.

When your API clients are ready, add credential information to your `.tf` configuration files. See [Authenticate the Akamai Terraform Provider](../guides/akamai_provider_auth.md)
for details on creating API clients and available authentication methods. Once you're done authenticating, come back here to complete the Akamai Provider setup.

**Note:** Depending on the contract and group you select, the Edge DNS and Global Traffic Management (GTM) modules may interact Property Manager API (PAPI). If so, be sure to include PAPI authorization in the API Clients for Edge DNS and GTM.

## Retrieve contract and group IDs

You'll need contract and group IDs to use most Akamai Provider modules.

You can retrieve these IDs through the [`akamai_contract`](../data-sources/property_contract.md) and
[`akamai_group`](../data-sources/property_group.md) data sources, which require read access to the Property Manager API. You can also get this information from the Contracts app in Akamai
Control Center, or by using other Akamai APIs or CLIs.

### Retrieve contract IDs with akamai_contract

You can get your contract ID automatically using the [`akamai_contract` data source](../data-sources/property_contract.md). This data source requires access to the Property Manager (PAPI) API service. See [Set up your API clients](../guides/akamai_provider_auth.md#set-up-your-api-clients).

To retrieve the default contract, you need to enter a group name or ID. No attributes need to be set:

```hcl
data "akamai_contract" "default" {
     group_name = "example group name"
}
```

You can now refer to the contract ID using the `id` attribute: `data.akamai_contract.default.id`.

### Retrieve group IDs with akamai_group

Akamai groups control access to your Akamai configurations and help consolidate reporting functions. Each account features a hierarchy of groups, which typically map to an organizational hierarchy.

Your account admins can use Control Center or the [Identity Management: User Administration API](https://developer.akamai.com/en-us/api/core_features/identity_management_user_admin/v2.html)
to set up groups, each with their own set of users and roles.

You can get your group ID automatically using the [`akamai_group` data source](../data-sources/property_group.md). To retrieve the default group ID you need to enter a contract ID:


```hcl
data "akamai_group" "default" {
	contract_id = data.akamai_contract.default.id
}
```

To get a specific group, you can specify the `name` argument:

```hcl
data "akamai_group" "default" {
	name = "example"
	contract_id = data.akamai_contract.default.id
}
```

You can now refer to the group ID using the `id` attribute: `data.akamai_group.default.id`.

## Arrange resources and data sources in the Akamai configuration file

You're now ready to import the existing configurations or create new ones from scratch.

At this point in the setup, you should refer to the guides for the Akamai modules you're using:

| **Module** | **Guide** |
|------------|------------|
| **Application Security** (beta) | [Application Security Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_appsec) |
| **Certificate Provisioning** | [Certificate Provisioning Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_cps) |
| **DNS Zone Administration** | [DNS Zone Administration Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_dns_zone) |
| **Global Traffic Management Domain Administration** | [Global Traffic Management Domain Administration Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_gtm_domain) |
| **Identity and Access Management** | [Identity and Access Management Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_iam) |
| **Network Lists** | [Network Lists Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_networklists) |
| **Property Provisioning** | [Property Provisioning Module Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_property) |

Once you're done with the module-level setup, continue with the next
sections here to initialize Akamai Provider, test the configuration, and apply the actions.

## Initialize the Akamai Provider

Once you have your configuration complete, save the `.tf` files. Then
switch to the terminal to initialize Terraform using the command:

    $ terraform init

This command installs the latest version of the Akamai Provider, as well
as any other providers you're using. To update
the Akamai Provider version after a new release, simply run `terraform
init` again.

## Test your configuration

To test your configuration, use Terraform's `plan` command:

    $ terraform plan

This command makes Terraform create a plan for the work it will do
based on the `.tf` configuration file. It doesn't actually make any changes
and is safe to run as many times as you like.

## Apply your configuration

You can execute all the actions you set in the configuration by running:

    $ terraform apply

The `apply` command previews all the changes before executing them, similarly to `plan`. Unless you set the `-auto-approve` flag, you need to confirm you want to proceed with the operation and propagate changes to the Akamai platform.
