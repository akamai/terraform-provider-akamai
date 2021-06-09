---
layout: "akamai"
page_title: "Akamai: SelectedHostnames"
subcategory: "Application Security"
description: |-
  SelectedHostnames
---

# akamai_appsec_selected_hostnames


The `akamai_appsec_selected_hostnames` resource allows you to set the list of hostnames protected under a given security configuration.


## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostnames = [ "example.com" ]
  mode = "APPEND"
}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `hostnames` - (Required) The list of hostnames to be applied, added or removed.

* `mode` - (Required) A string specifying the interpretation of the `hostnames` parameter. Must be one of the following:

  * APPEND - the hosts listed in `hostnames` will be added to the current list of selected hostnames
  * REPLACE - the hosts listed in `hostnames` will overwrite the current list of selected hostnames
  * REMOVE - the hosts listed in `hostnames` will be removed from the current list of select hostnames

# Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

