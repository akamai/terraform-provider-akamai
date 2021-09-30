---
layout: "akamai"
page_title: "Module: DNS Zone Administration"
description: |-
  DNS Zone Administration module for the Akamai Terraform Provider
---

# DNS Zone Administration Guide

The Akamai Provider for Terraform provides you the ability to automate the creation, deployment, and management of DNS zone configuration and administration; as well as 
importing existing zones and recordsets.  

To get more information about Edge DNS, see:

* Developer - [API documentation](https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html).
* User Guide - [Official Documentation](https://learn.akamai.com/en-us/products/cloud_security/edge_dns.html).


## Prerequisites 

Before starting with the DNS module, you need to:

1. Complete the tasks in [Get Started with the Akamai Provider](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_provider). 
2. Determine whether you want to import an existing DNS zone and records or create new ones.
3. If you're importing an existing DNS configuration, continue with [Import a DNS zone and records](#import-a-dns-zone-and-records).
3. If you're creating a new DNS configuration, continue with [Create a DNS zone
](#create-a-dns-zone) and [Create a DNS record](#create-a-dns-record).

## Import a DNS zone and records

You can migrate an existing Edge DNS zone into your Terraform configuration using either a command line utility or step-by-step construction.

### Import using the command line utility

You can use the [Akamai CLI for Akamai Terraform Provider](https://github.com/akamai/cli-terraform) to generate a configuration for and import an existing Edge DNS zone and its recordsets. With the package, you can generate:

* a JSON-formatted list of the zone and recordsets.
* a Terraform configuration for the zone and select recordsets.
* a command line script to import all defined resources.

Before using this CLI, keep the following in mind:

* Download the existing zone configuration and master file to have as a backup and reference during an import. You can download these by using the [Edge DNS Zone Management API](https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html) or the Edge DNS app on [Control Center](https://control.akamai.com).  
* Terraform limits the characters that can be part of its resource names. During construction of the resource configurations, invalid characters are replaced with underscore , '_'.
* Terraform doesn't provide any state information during import. When you run `plan` and `apply` after an import, Terraform lists discrepencies and reconciles configurations and state. Any discrepencies clear following the first `apply`. 
* After first time you run `plan` or `apply`, the `contract` and `group` attributes are updated.
* Run `terraform plan` after importing to validate the generated `tfstate` file.

### Import using step-by-step construction

To import using step-by-step construction, complete these tasks:

1. Determine how you want to test your Terraform import. For example, you may want to set up your zone and recordset imports in a test environment to familiarize yourself with the provider operation and mitigate any risks to your existing DNS zone configuration.
1. Download the existing zone configuration and master file to have as a backup and reference during an import. You can download these from the [Edge DNS Zone Management API](https://developer.akamai.com/api/cloud_security/edge_dns_zone_management/v2.html) or from the Edge DNS app on [Control Center](https://control.akamai.com) .  
1. Using the zone master file as a reference, create a Terraform configuration representing the existing zone and all contained recordsets. 
1. Verify that your Terraform configuration addresses all required attributes and any optional and computed attributes you need.
1. Run `terraform import`. This command imports the existing zone and contained recordsets. The import happens in serial order.
1. Compare the downloaded zone master file with the `terraform.tfstate` file to confirm that the zone and all recordsets are represented correctly.
1. Run `terraform plan` on the configuration. The plan should be empty. If not, correct accordingly and repeat until plan is empty and configuration is in sync with the Edge DNS backend.

## Create a DNS zone

The zone itself is represented by a [`akamai_dns_zone` resource](../resources/dns_zone.md). Add this new resource block to your `akamai.tf` file after the provider block. **Note:** the zone should be the first DNS resource created as it provides operating context for all other recordset resources.

To define the entire configuration, we start by opening the resource block and giving the `zone` a name. In this case we're going to use the name "example."

Next, we set the required (`zone`, `type`, `group`, `contract`) and optional (`comment`) arguments for a simpler secondary `type`.

Once done, your `akamai.tf` configuration file should include configuration items such as:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
      version = "1.2.0"
    }
  }
}

locals {
	section = "default"
}

provider "akamai" {
	edgerc = "~/.edgerc"
	config_section = local.section
}

data "akamai_contract" "default" { }

data "akamai_group" "default" {
	contract_id = data.akamai_contract.default.id
}

resource "akamai_dns_zone" "example_com" {
	zone = "examplezone.com"                      # Zone Name
	type = "secondary"				              # Zone type
	masters = [ "1.2.3.4" ]				          # Zone master(s)
	group = data.akamai_group.default.id          # Group ID variable
	contract = data.akamai_contract.default.id    # Contract ID variable
	comment = "example zone demo"
}
```
> **Note:** Notice the use of variables from the previous section to reference the group and contract IDs. These will be replaced at runtime by Terraform with the actual values.

### Validate Terraform Zone Configuration and State

To validate the configuration up to this point, run the following command. The actual commit will come later in the procedure with an apply command.

```
$ terraform plan
```

### Primary Zones

Unlike creating secondary zone types, creating primary zone types is best by following a multi-step process as follows. To complete these steps, you need to download and install the [Akamai CLI](https://developer.akamai.com/cli) and [CLI-Terraform package](https://github.com/akamai/cli-terraform). 

#### Configure Zone

In addition to `akamai.tf` set with Get Started with the [Akamai Terraform Provider Guide](https://registry.terraform.io/providers/akamai/akamai/latest/docs/guides/get_started_provider), create the zone configuration in a new zone configuration file. For this example, use `example_primary_zone_com.tf`.

**Note:** Subsequent steps will require the zone configuration file be named `<zone>.tf` with dots replaced by underscores. Edge DNS will automatically create NS and SOA records. Steps below show how to synchronize these records to the local Terraform state.

##### Example configuration:

```
locals {
	section     = "default"
	zone        = "example_primary_zone.com"
}

provider "akamai" {
	edgerc = "~/.edgerc"
	config_section = local.section
}

data "akamai_contract" "default" { 
    group_name = "Example group"
}

data "akamai_group" "default" {
	contract_id = data.akamai_contract.default.id
}

resource "akamai_dns_zone" "primary_example" {
	zone = local.zone
	type = "primary"
	group    = data.akamai_group.default.id
	contract = data.akamai_contract.default.id
	comment = "example primary zone and records"
}
```

**Note:** Referencing items in the locals block is done so with a singular `local` prefix such as `local.section`. Because Terraform references variables in all `.tf` files, the locals and provider blocks may not necessary in this zone file.

### Validate Terraform Zone Configuration and State

To validate the configuration up to this point, run the following command. The actual commit will come later in the procedure with an apply command.

```
$ terraform plan
```

**Note:** You can run `terraform plan` many times.

### Adding Zone SOA and NS Records To TF Configuration

Creating a primary zone has the side effect of creating both initial SOA and NS records. Without these two recordsets, the zone cannot be managed. Using the CLI-Terraform CLI package, the zone's top level SOA and NS records now need to be added to the Terraform configuration as follows.

#### Create a List of Zone Recordsets

First, create a list of the zone's current recordsets.

```
$ akamai terraform create-zone example_primary_zone.com --resources
```

The command will generate a file, `example_primary_zone_com_resources.json`, with the following content:

```
{
  "Zone": "example_primary_zone.com",
  "Recordsets": {
    "example_primary_zone.com": [
      "NS",
      "SOA"
    ]
  }
}
```

#### Update the Terraform Zone Configuration File

Next, update the Terraform Zone configuration file using the previously generated JSON file as input and the following command.

```
$ akamai terraform create-zone example_primary_zone.com --createconfig
```

The zone configuration file, `example_primary_zone_com.tf`, will be updated with the resulting content:

```
resource "akamai_dns_zone" "primary_example" {
	zone = "local.zone
	type = "primary"
	group    = data.akamai_group.default.id
	contract = data.akamai_contract.default.id
	comment = "example primary zone and records"
}

resource "akamai_dns_record" "example_primary_zone_com_example_primary_zone_com_NS" {
	zone = local.zone
	recordtype = "NS"
	ttl = 86400
	target = ["ax-xx.akam.net.", "axx-xx.akam.net.", "axx-xx.akam.net.", "ax-xx.akam.net.", "ax-xx.akam.net.", "ax-xx.akam.net."]
	name = "example_primary_zone.com"
}

resource "akamai_dns_record" "example_primary_zone_com_example_primary_zone_com_SOA" {
	zone = local.zone
	expiry = 604800
	nxdomain_ttl = 300
	name = "example_primary_zone.com"
	target = []
	name_server = "ax-xx.akam.net."
	email_address = "hostmaster.example_primary_zone.com."
	refresh = 3600
	retry = 600
	recordtype = "SOA"
	ttl = 86400
}
```
**Note:** Name server targets have been masked. Also, a default `dnsvars.tf` file is generated. It can be ignored, deleted or used. Other Terraform configuration files can reference variables in this file with a macro such as "${dnsvar.zone}".

#### Generate a Resource Import Script

Next, generate a zone resources import script using previously generated output.

```
$ akamai terraform create-zone example_primary_zone.com --importscript
```

The file `example_primary_zone.com_resource_import.script` is generated with the following content:

```
terraform init
terraform import akamai_dns_zone.clidns_primary_test_com clidns_primary_test.com
terraform import akamai_dns_record.clidns_primary_test_com_clidns_primary_test_com_NS clidns_primary_test.com#clidns_primary_test.com#NS
terraform import akamai_dns_record.clidns_primary_test_com_clidns_primary_test_com_SOA clidns_primary_test.com#clidns_primary_test.com#SOA
```

Next, edit the script file and remove the line `terraform import akamai_dns_zone.egl_clidns_primary_test_com egl_clidns_primary_test.com` as the zone does not need to be imported.

#### Import Zone Recordsets

Perform the following command to import the recordsets into Terraform.

```
$ ./example_primary_zone.com_resource_import.script
```

The Terraform configuration and state will now contain the zone's SOA and NS Records with values consistent with the Akamai DNS Infrastructure.

### Validate Terraform Zone Configuration and State

To validate the configuration up to this point, run the following command. The actual commit will come later in the procedure with an apply command.

```
$ terraform plan
```

## Create a DNS record

The recordset itself is represented by a [`akamai_dns_record` resource](../resources/dns_record.md). Add this new block to your `akamai.tf` file after the provider block.

To define the entire configuration, we start by opening the resource block and give it a name. In this case we're going to use the name "example_a_record".

Next, we set the required (zone, recordtype, ttl) and any optional/required arguments based on recordtype. Required fields for each record type are itemized in [`akamai_dns_record` resource](../resources/dns_record.md).

Once complete, your record configuration should look like this:

```
resource "akamai_dns_record" "example_a_record" {
	zone = akamai_dns_zone.example.zone
	target = ["10.0.0.2"]
	name = "example_a_record"
	recordtype = "A"
	ttl = 3600
}
```

## Validate Terraform Zone Configuration and State

To validate the configuration up to this point, run the following command. The actual commit will come later in the procedure with an apply command.

```
$ terraform plan
```

## Apply Changes

To actually create our zone and recordset, we need to instruct Terraform to apply the changes outlined in the plan. To do this, run the command:

```
$ terraform apply
```

Once this completes, your zone and recordset will have been created. You can verify this in [Akamai Control Center](https://control.akamai.com).

## Import Records

Existing DNS resources may be imported using one of the following formats:

```
$ terraform import akamai_dns_zone.{{zone resource name}} {{edge dns zone name}}
$ terraform import akamai_dns_record.{{record resource name}} {{edge dns zone name}}#{{edge dns recordset name}}#{{record type}}
```

## How you can use the DNS module

These sections include information on different ways to use the Akamai's Terraform DNS module:

* [Working With MX Records](#working-with-mx-records)
* [Important Behavior Considerations](#important-behavior-considerations)
* [Primary Zone Partially Created](#primary-zone-partially-created)

## Working With MX Records

MX Record resource configurations may be instantiated in three different forms:

1. Coupling Priority and Host.
2. Assigning Priority to Hosts via Variables.
3. Instance Generation.

### Coupling Priority and Host

With this configuration style, each target entry includes both the priority and host. The following configuration will produce a recordset rdata value of:

```
["0 smtp-0.example.com.", "10 smtp-1.example.com."]
```

```
resource "akamai_dns_record" "mx_record_self_contained" {
	zone = local.zone
	target = ["0 smtp-0.example.com.", "10 smtp-1.example.com."]
	name = "mx_record_self_contained.example.com"
	recordtype = "MX"
	ttl = 300
}
```

### Assigning Priority to Hosts via Variables

With this configuration style, a number of hosts will be defined in the target field as a list. A starting `priority` and `priority_increment` are also defined. The provider 
will construct the rdata values by incrementally pairing and incrementing the `priority` by the `priority_increment`. For example, the following configuration will produce a recordset rdata value of:

```
["10 smtp-1.example.com.", "20 smtp-2.example.com.", "30 smtp-3.example.com."]
```

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

### Instance Generation

With this configuration style, a number of host instances can be generated using Terraform's count or for/each construct. For example, the following configuration will produce three distinct resource instances, each with a single target and priority, and an aggregated recordset rdata value of: 

```
["0 smtp-0.example.com.", "10 smtp-1.example.com.", "20 smtp-2.example.com."] 
```

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

## Important Behavior Considerations

* Concurrent and independent modifications through the Terraform provider and Control Center UI may result in configuration drift and require manual intervention to reconcile the local Terraform state. This issue is particularly a concern for MX records.
* Deletion of a record resource with multiple instances or deletion of a single instance, will result in the entire remote recordset resource being removed.
* Record configurations and state include a computed `record_sha` field that represents the current resource state to compare the local and remote MX record configurations. This field will not exist in upgraded configurations. As such, doing a plan on an existing MX record may result in the following message to ignore.

```
No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.
```

## Primary Zone Partially Created

While it's rare, sometimes a primary zone is only partially created on the Akamai backend. For example, a network error happens and while the zone was created, the SOA and NS records were not. 

Any attempt to manage or administer recordsets in the zone will fail. To resolve this issue, you have to manually create the SOA and NS records before you can manage the configuration.

You can create these records from the Edge DNS application available from [Akamai Control Center](https://control.akamai.com). You also have the option of using the [Akamai CLI for Edge DNS](https://github.com/akamai/cli-dns).
