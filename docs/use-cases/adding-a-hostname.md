---
layout: "akamai"
page_title: "Adding a Hostname to a Security Configuration"
description: |-
  Adding a Hostname to a Security Configuration
---


# Adding a Hostname to a Security Configuration

Security configurations are not designed to automatically protect every server in your organization; instead, each security configuration protects only those servers (i.e., only those hosts) you've designated for this protection. Do you have different servers that require different levels or different types of protection? That's fine: put one set of servers in Configuration A, another set in Configuration B, and so on. This gives you the ability to fine-tune your websites and to strike a balance between security and efficiency.

So how do you designate a host for protection by a security configuration? That turns out to be a two-step process:

1.	Determine which hosts are “selectable;” that is, which hosts are available to be added to a security configuration. Note that it's possible that some of your hosts *can't* be added to a security configuration. For example, if a contract has expired you might not be able to offer protection to certain servers. Note also that a single security configuration can manage multiple hosts. However, an individual host can only be a member of a single security configuration. In other words, Host A can belong to Configuration A or Configuration B, but it can't belong to both Configuration A and Configuration B at the same time.

2.	Add the hostname to the appropriate security configuration. When updating the selected hosts for a security configuration, you can either append the new host (or hosts: you can add multiple hosts in a single operation) to the existing collection of hosts, or you can replace the existing collection with the hosts specified in your Terraform configuration.

When adding multiple hosts you can specify each host individually, or (if possible) you can use wildcard characters. For example, this syntax adds four hosts (**host1**, **host2**, **host3**, and **host4**) from the akamai.com domain to a security configuration:

```
hostnames = ["host1.akamai.com",  "host2.akamai.com", "host3.akamai.com", "host4.akamai.com"]
```

That syntax works just fine. But so does this syntax, which adds all of your `akamai.com` hosts to the configuration:

```
hostnames = [ "*.akamai.com"]
```

## Returning a Collection of Selectable Hostnames

To return a list of your available hostnames, use the [akamai_appsec_selectable_hostname](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_selectable_hostnames)s data source. There are two different ways to use this data source: you can return all the available hostnames for a specific contract and group, or you can return all the available hostnames for a specific configuration. For example, this Terraform configuration returns all the selectable hostnames for contract **7-5DR99** and group **57112**:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_selectable_hostnames" "selectable_hostnames_for_create_configuration" {
  contractid = "7-5DR99"
  groupid    = 57112
}

output "selectable_hostnames_for_create_configuration" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames_for_create_configuration.hostnames
}
```

Alternatively, this configuration returns the selectable hostnames for configuration **58843**:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_selectable_hostnames" "selectable_hostnames" {
  config_id = 58843
}

output "selectable_hostnames" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames
}
```

Either way, if you call this configuration by using the `terraform plan` command, you'll get back a list of hostnames similar to the following

```
Changes to Outputs:

selectable_hostnames_for_create_configuration = [
  "host1.akamai.com",
  "host2.akamai.com",
  "host3.akamai.com",
  "host4.akamai.com"
]
```

Any (or all) of these hostnames can be added to your configuration.

Incidentally, if you try adding a hostname that *isn't* on the list of selectable hostnames then your command will fail with an error similar to this:

```
╷
│ Error: setting property value: Title: Invalid Input Error; Type: https://problems.luna.akamaiapis.net/appsec-configuration/error-types/INVALID-INPUT-ERROR; Detail: You don't have access to add these hostnames [identitydocs.akamai.com]
│
│   with akamai_appsec_selected_hostnames.appsecselectedhostnames,
│   on akamai.tf line 15, in resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames":
│   15: resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
```

## Adding a Hostname to a Security Configuration

After you've chosen a hostname (or set of hostnames), you can add these hosts to a security configuration by using the [akamai_appsec_selected_hostnames](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_selected_hostnames) resource and a Terraform configuration similar to this:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostnames = [ "host1.akamai.com", "host2.akamai.com"]
  mode      = "APPEND"
}
```

In this configuration, we start by:

1.	Declaring the Akamai provider.
2.	Providing our authentication credentials.
3.	Connecting to the **Documentation** security configuration.

When all that's done, we use this block of code the add the hostnames **host1.akamai.com** and **host2.akamai.com** to the **Documentation** configuration:

```
resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostnames = [ "host1.akamai.com", "host2.akamai.com"]
  mode = "APPEND"
}
```

Note that, in order to add the hostnames to the existing collection of hostnames for the configuration, we set the mode to `APPEND`. For example, suppose our collections currently has these three hosts:

- hostA.akamai.com
- hostB.akamai.com
- hostC.akamai.com

After we run our Terraform configuration the security configuration will have these hosts:

- hostA.akamai.com
- hostB.akamai.com
- hostC.akamai.com
- host1.akamai.com
- host2.akamai.com

By comparison, if we set mode to `REPLACE`, the existing set of hostnames will be deleted and will be *replaced* by the hostnames specified in the Terraform configuration. In other words, the security configuration will end up with these hosts:

- host-1.akamai.com
- host-2.akamai.com
