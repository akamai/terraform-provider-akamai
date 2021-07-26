---
layout: "akamai"
page_title: "Akamai: WAPSelectedHostnames"
subcategory: "Application Security"
description: |-
  WAPSelectedHostnames
---

# akamai_appsec_wap_selected_hostnames [Beta]


The `akamai_appsec_wap_selected_hostnames` resource allows you to set the lists of hostnames to be protected and to be evaluated
under a given security configuration and policy. This resource is available only for WAP accounts. Either of the lists of hostnames
may be omitted or specified as an empty list, but at least one of the two lists must be present and non-empty. (WAP selected hostnames
is currently in beta. Please contact your Akamai representative for more information about this feature.)

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

resource "akamai_appsec_wap_selected_hostnames" "appsecwap_selectedhostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  protected_hostnames = [ "example.com" ]
  evaluated_hostnames = [ "example2.com" ]
}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `protected_hostnames` - (Optional) The list of hostnames to be protected. If not supplied, then `evaluated_hostnames` must be supplied.

* `evaluated_hostnames` - (Optional) The list of hostnames to be evaluated. If not supplied, then `protected_hostnames` must be supplied.

# Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

