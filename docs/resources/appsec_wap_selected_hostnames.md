---
layout: "akamai"
page_title: "Akamai: WAPSelectedHostnames"
subcategory: "Application Security"
description: |-
  WAPSelectedHostnames
---

# akamai_appsec_wap_selected_hostnames


The `akamai_appsec_wap_selected_hostnames` resource allows you to set the lists of hostnames to be protected and to be evaluated
under a given security configuration and policy. This resource is available only for WAP accounts.

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

* `protected_hostnames` - (Required) The list of hostnames to be protected.

* `evaluated_hostnames` - (Required) The list of hostnames to be evaluated.

# Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

