---
layout: "akamai"
page_title: "Akamai: SelectedHostname"
subcategory: "APPSEC"
description: |-
  SelectedHostname
---

# resource_akamai_appsec_selected_hostname


The `resource_akamai_appsec_selected_hostname` resource allows you to create or re-use SelectedHostnames.

If the SelectedHostname already exists it will be used instead of creating a new one.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Example for EDGE"
  
}



output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}

resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version 
    hostnames = ["rinaldi.sandbox.akamaideveloper.com","sujala.sandbox.akamaideveloper.com"]  
}

```

## Argument Reference

The following arguments are supported:
* `config_id`- (Required) The Configuration ID

* `version` - (Required) The Version Number of configuration

* `hostnames` - (Required) The List of hostnames to configure

# Attributes Reference

The following are the return attributes:

* `Hostnames` - Set of selectable hostnames

