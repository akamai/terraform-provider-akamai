---
layout: "akamai"
page_title: "Akamai: Get Started with Network Lists"
description: |-
  Get Started with Akamai Network Lists using Terraform
---

# Get Started with Network Lists

The Akamai Network Lists provider for Terraform gives you the ability to automate the creation, deployment, and management of lists used in various Akamai security products such as Kona Site Defender, Web App Protector, and Bot Manager. Network lists are shared sets of IP addresses, CIDR blocks, or broad geographic areas. Along with managing your own lists, you can also access read-only lists that Akamai dynamically updates for you. For more information about Network Lists, see the [API documentation](https://developer.akamai.com/api/cloud_security/network_lists/v2.html)

## Configure the Terraform Provider

Set up your .edgerc credential files as described in [Get Started with Akamai APIs](https://developer.akamai.com/api/getting-started), and include read-write permissions for the Network Lists API. 

1. Create a new folder called `terraform`
1. Inside the new folder, create a new file called `akamai.tf`.
1. Add the provider configuration to your `akamai.tf` file:

```hcl
provider "akamai" {
	edgerc = "~/.edgerc"
	config_section = "networklists"
}
```

## Prerequisites

Review [Get Started with APIs](https://learn.akamai.com/en-us/learn_akamai/getting_started_with_akamai_developers/developer_tools/getstartedapis.html) for details on how to set up client tokens to access any Akamai API. These tokens appear as custom hostnames that look like this: https://akzz-XXXXXXXXXXXXXXXX-XXXXXXXXXXXXXXXX.luna.akamaiapis.net.

To enable this API, choose the API service named Network Lists, and set the access level to READ-WRITE.

## Retrieving Network List Information

You can obtain a list of all network lists available for an authenticated user belonging to a group using the [`akamai_networklists_network_lists`](../data-sources/networklists_network_lists.md) data source and its `output_text` attribute. Add the following to your `akamai.tf` file:

```hcl
data "akamai_networklist_network_lists" "network_lists" {
}

output "network_lists_text" {
  value = data.akamai_networklist_network_lists.network_lists.output_text
}
```

Once you have saved the file, switch to the terminal and initialize Terraform using the command:

```bash
$ terraform init
```

This command will install the latest version of the Akamai provider, as well as any other providers necessary. To update the Akamai provider version after a new release, simply run `terraform init` again.

## Test Your Configuration

To test your configuration, use `terraform plan`:

```bash
$ terraform plan
```

This command will make Terraform create a plan for the work it will do based on the configuration file. This will not actually make any changes and is safe to run as many times as you like.

## Apply Changes

To actually display the configuration information, or to create or modify resources as described further in this guide, we need to instruct Terraform to `apply` the changes outlined in the plan. To do this, in the terminal, run the command:

```bash
$ terraform apply
```

Once this command has been executed, Terraform will display to the terminal window a formatted display of the ID, name, type, elementCount, syncPoint and readonly status of the existing network lists. The `json` attribute of the `networklist_network_lists` data source will produce a JSON-formatted output containing similar information.

You can filter the network list output by supplying additional parameters to the `networklist_network_lists` data source. The `name` and `type` parameters will limit the output to the list with the specified values. Add the following example of filtering to your `config.tf` file:

```hcl
data "akamai_networklist_network_lists" "network_lists_filter" {
  name = "test-network-list1"
  type = "IP"
}
```

## Create a Network List

You can create a new network list using the `resource_networklist_network_list` resource. Unlike Terraform data sources, a Terraform resource is capable of making changes to your configuration by creating or modifying objects. The `resource_networklist_network_list` resource requires parameters that indicate the type of list (`IP` or `GEO`) and the specific items to be included in the list (either IP addresses or locations), as well as the name of the list and a description. Lastly, the `mode` attribute indicates whether the items in the `list` parameter are to be added or removed from the indicated list, or replace the list contents entirely. Thus, you can create a new list or add to an existing list by specifying `append`, or remove elements from a network list by specifying `remove`, or replace the contents of an existing list entirely by specifying `replace`. Create a new network list by adding the following to your `config.tf` file.

```hcl
resource "akamai_networklist_network_list" "network_list" {
  name = "Test-whitelist-NL"
  type = "IP"
  description = "Network List description"
  list = [
    "13.230.0.0/15",
    "195.7.50.194",
    "50.23.59.233"
  ]
  mode = "APPEND"
}
```

Test your configuration by running `terraform plan` as above. You should see a formatted description of the network lists that will be created. To cause these changes to take effect, run `terraform apply`.

## Activate a Network List

You can activate a network list by using the `akamai_networklist_activation` resource. Add the following to your `config.tf` file:

```hcl
resource "akamai_networklist_activation" "activation" {
  network_list_id = data.akamai_networklist_network_lists.network_lists_filter.list[0]
  network = "STAGING"
  notes  = "TEST Notes"
  notification_emails = ["user@example.com"]
}
```

This example uses the ID of the first element in the `network_lists_filter` example seen earlier. The Terraform provider activates this network list and checks on its progress as the activation proceeds. Once the operation is complete, the provider generates an email to each of the addresses in the `notification_emails` list.

## Subscribe to a Network List

You can subscribe one or more email addresses to receive notifications when any of a set of network lists are modified. Add the the following to your `config.tf` file:

```
resource "akamai_networklist_subscription" "subscribe" {
  network_list = data.akamai_networklist_network_lists.network_lists_filter.list
  recipients = ["user@example.com"]
}
```
Once you `apply` these changes, the `user@example.com` address will be notified when any of the lists in the `network_lists_filter.list` set are modified.

