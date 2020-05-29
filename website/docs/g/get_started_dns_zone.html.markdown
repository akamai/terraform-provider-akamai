---
layout: "akamai"
page_title: "Akamai: Get Started with DNS Zone Administration"
sidebar_current: "docs-akamai-guide-get-started-dns-zone"
description: |-
  Get Started with Akamai DNS Zone Administration using Terraform
---

# Get Started with DNS Zone Administration

The Akamai Provider for Terraform provides you the ability to automate the creation, deployment, and management of DNS zone configuration and administration; as well as 
importing existing zones and recordsets.  

To get more information about Edge DNS, see:

* [API documentation](https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html)
* How-to Guides
    * [Official Documentation](https://learn.akamai.com/en-us/products/cloud_security/edge_dns.html)

## Configure the Terraform Provider

Set up your credential files as described in [Get Started with Akamai APIs](https://developer.akamai.com/api/getting-started), and include authorization for the GTM Config API 

Next, we need to configure the provider with our credentials. This is done using a provider configuration block.

1. Create a new folder called `terraform`
1. Inside the new folder, create a new file called `akamai.tf`.
1. Add the provider configuration to your `akamai.tf` file:

```hcl
provider "akamai" {
    dns {
        host = "..."
        access_token = "..."
        client_token = "..."
        client_secret = "..."
    }
}
```

## Prerequisites

To create a zone there are several dependencies you must first meet:

* **Contract ID**: The ID of the contract under which the zone and contained recordsets will live
* **Group ID**: The ID of the group under which the zone and contained recordsets will live

To import an existing zone and recordsets, you must also know the identifiers or the objects; e.g. zone and recordset names in addition to the prior information.

## Retrieving The Contract ID

You can fetch your contract ID automatically using the [`akamai_contract` data source](/docs/providers/akamai/d/contract.html). To fetch the default contract ID no attributes need to be set:

```hcl
data "akamai_contract" "default" {

}
```

Alternatively, if you have multiple contracts, you can specify the `group` which contains it:

```hcl
data "akamai_contract" "default" {
  group = "default"
}
```

You can now refer to the contract ID using the `id` attribute: `data.akamai_contract.default.id`.

## Retrieving The Group ID

Similarly, you can fetch your group ID automatically using the [`akamai_group` data source](/docs/providers/akamai/d/group.html). To fetch the default group ID no attributes need to be set:

```hcl
data "akamai_group" "default" {

}
``` 

To fetch a specific group, you can specify the `name` argument:

```hcl
data "akamai_group" "default" {
  name = "example"
}
```

You can now refer to the group ID using the `id` attribute: `data.akamai_group.default.id`.

## Creating a DNS Zone

The zone itself is represented by a [`akamai_dns_zone` resource](/docs/providers/akamai/r/dns_zone.html). Add this new resource block to your `akamai.tf` file after the provider block. Note: the zone should be the first DNS resource created as it provides operating context for all other recordset resources.

To define the entire configuration, we start by opening the resource block and giving the zone a name. In this case we’re going to use the name "example".

Next, we set the required (zone, type, group, contract) and optional (comment) arguments.

Once you’re done, your zone configuration should look like this:

```hcl
resource "akamai_dns_zone" "example" {
        zone = "examplezone.com"                        # Zone Name
        type = "primary"				# Zone type
        group    = data.akamai_group.default.id         # Group ID variable
        contract = data.akamai_contract.default.id      # Contract ID variable
	comment = "example zone demo"
}
```
> **Note:** Notice that we’re using variables from the previous section to reference the group and contract IDs. These will automatically be replaced at runtime by Terraform with the actual values.

## Creating a DNS Record

The recordset itself is represented by a [`akamai_dns_record` resource](/docs/providers/akamai/r/dns_record.html). Add this new block to your `akamai.tf` file after the provider block.

To define the entire configuration, we start by opening the resource block and give it a name. In this case we’re going to use the name "example_a_record".

Next, we set the required (zone, recordtype, ttl) and any optional/required arguments based on recordtype. Required fields for each record type are itemized in [`akamai_dns_record` resource](/docs/providers/akamai/r/dns_record.html).

Once you’re done, your record configuration should look like this:

```hcl
resource "akamai_dns_record" "example_a_record" {
    zone = akamai_dns_zone.example.zone
    target = ["10.0.0.2"]
    name = "example_a_record"
    recordtype = "A"
    ttl = 3600
}
```

## Initialize the Provider

Once you have your configuration complete, save the file. Then switch to the terminal to initialize terraform using the command:

```bash
$ terraform init
```

This command will install the latest version of the Akamai provider, as well as any other providers necessary (such as the local provider). To update the Akamai provider version after a new release, simply run `terraform init` again.

## Test Your Configuration

To test your configuration, use `terraform plan`:

```bash
$ terraform plan
```

This command will make Terraform create a plan for the work it will do based on the configuration file. This will not actually make any changes and is safe to run as many times as you like.

## Apply Changes

To actually create our zone and recordset, we need to instruct terraform to apply the changes outlined in the plan. To do this, in the terminal, run the command:

```bash
$ terraform apply
```

Once this completes your zone and recordset will have been created. You can verify this in [Akamai Control Center](https://control.akamai.com).

## Import

Existing DNS resources may be imported using one of the following formats:

```
$ terraform import akamai_dns_zone.{{zone resource name}} {{edge dns zone name}}
$ terraform import akamai_dns_record.{{record resource name}} {{edge dns zone name}}#{{edge dns recordset name}}#{{record type}}
```

[Migrating A DNS Zone](/docs/providers/akamai/g/faq.html#migrating-an-edge-dns-zone-and-records-to-terraform) discusses DNS resource import in more detail.

## Working With MX Records

MX Record resource configurations may be instantiated in a number of different forms. These forms consist of:

### Coupling Priority an Host

With this configuration style, each target entry includes both the priority and host. The following configuration

```
resource "akamai_dns_record" "mx_record_self_contained" {
    zone = local.zone
    target = ["0 smtp-0.example.com.", "10 smtp-1.example.com."]
    name = "mx_record_self_contained.example.com"
    recordtype = "MX"
    ttl = 300
}
```
will produce a recordset rdata value of 

```
["0 smtp-0.example.com.", "10 smtp-1.example.com."]
```

### Assigning Priority to Hosts via Variables

With this configuration style, a number of hosts will be defined in the target field as a list. A starting priority and priority_increment are also defined. The provider 
will construct the rdata values by incrementally pairing and incrementing the priority by the priority_increment. For example, the following configuration

```
resource "akamai_dns_record" "mx_record_pri_increment" {
    zone = local.zone
    target = ["smtp-1.example.com.", "smtp-2.example.com.", "smtp-3.example.com."]
    priority = 10
    priority_increment = 10
    name = "mx_pri_increment.example.com"
    recordtype = "MX"
    ttl = 900
}
```
will produce a recordset rdata value of 

```
["10 smtp-1.example.com.", "20 smtp-2.example.com.", "30 smtp-3.example.com."]
```

### Instance Generation

With this configuration style, a number of host instances can be generated by using Terraform's count or for/each construct. For example, the following configuration

```
resource "akamai_dns_record" "mx_record_instances" {
    zone = local.zone
    name = "mx_record.example.com"
    recordtype =  "MX"
    ttl =  500
    count = 3
    target = ["smtp-${count.index}.example.com."]
    priority = count.index*10
}
```
will produce three distinct resource instances, each with a single target and priority, and an aggregated recordset rdata value of 

```
["0 smtp-0.example.com.", "10 smtp-1.example.com.", "20 smtp-2.example.com."] 
```

## Important Behavior Considerations

* Concurrrent modifications thru the Terraform provider and the UI may result in configuration drift and require manual intervention to reconcile. This issue is particularly a concern for MX records.
* Deletion of a record resource with multiple instances or deletion of a single instance, will result in the entire remote recordset resource being removed.
* Record configurations and state include a computed record_sha field. This field is used represent the current resource state as well as to compare local MX record configuration with the remote configuration. This field will not exist in upgraded configurations. As such, doing a plan on an existing MX record may result in the  following informational message which can be ignored.

```
No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.
```
