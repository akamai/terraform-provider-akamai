---
layout: "akamai"
page_title: "Akamai: SelectedHostnames"
subcategory: "APPSEC"
description: |-
 SelectedHostnames
---

# akamai_appsec_selected_hostnames

Use `akamai_appsec_selected_hostnames` data source to retrieve a selected_hostnames id.

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

data "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version  
}


```

## Argument Reference

The following arguments are supported:

* `config_id`- (Required) The Configuration ID

* `version` - (Required) The Version Number of configuration

# Attributes Reference

The following are the return attributes:

* `Hostnames` - Set of selected hostnames

* `Hostnames_json` - Set of selected hostnames in json format

