---
layout: "akamai"
page_title: "Akamai: WAPSelectedHostnames"
subcategory: "Application Security"
description: |-
 WAPSelectedHostnames
---

# akamai_appsec_wap_selected_hostnames [Beta]

Use the `akamai_appsec_wap_selected_hostnames` data source to retrieve lists of the hostnames that are currently
protected and currently being evaluated under a given security configuration and policy. This resource is available
only for WAP accounts. (WAP selected hostnames is currently in beta. Please contact your Akamai representative for
more information about this feature.)

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

data "akamai_appsec_wap_selected_hostnames" "wap_selected_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}

output "protected_hostnames" {
  value = data.akamai_appsec_wap_selected_hostnames.wap_selected_hostnames.protected_hostnames
}

output "evaluated_hostnames" {
  value = data.akamai_appsec_wap_selected_hostnames.wap_selected_hostnames.evaluated_hostnames
}

output "json" {
  value = data.akamai_appsec_wap_selected_hostnames.wap_selected_hostnames.json
}

output "output_text" {
  value = data.akamai_appsec_wap_selected_hostnames.wap_selected_hostnames.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `protected_hostnames` - The list of hostnames currently protected under the given security configuration and policy.

* `evaluated_hostnames` - The list of hostnames currently being evaluated under the given security configuration and policy.

* `hostnames_json` - A JSON-formatted display of the protected and evaluated hostnames.

* `output_text` - A tabular display of the protected and evaluated hostnames.

