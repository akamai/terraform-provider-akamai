---
layout: "akamai"
page_title: "Akamai: SelectableHostnames"
subcategory: "APPSEC"
description: |-
 SelectableHostnames
---

# akamai_appsec_selectable_hostnames

Use `akamai_appsec_selectable_hostnames` data source to retrieve a selectable_hostnames id.

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

data "akamai_appsec_selectable_hostnames" "appsecselectablehostnames" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version   
}

output "selectablehostnames" {
  value = data.akamai_appsec_selectable_hostnames.appsecselectablehostnames.hostnames
}

```

## Argument Reference

The following arguments are supported:

* `config_id`- (Required) The Configuration ID

* `version` - (Required) The Version Number of configuration

* `active_in_staging` - (Optional) Active in staging

* `active_in_production` - (Optional) Active in production


# Attributes Reference

The following are the return attributes:

* `Hostnames` - Set of selectable hostnames

* `Hostnames_json` - Set of selectable hostnames in json format

