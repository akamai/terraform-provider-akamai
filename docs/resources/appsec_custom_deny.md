---
layout: "akamai"
page_title: "Akamai: CustomDeny"
subcategory: "Application Security"
description: |-
  CustomDeny
---

# resource_akamai_appsec_custom_deny

The `resource_akamai_appsec_custom_deny` resource allows you to create a new custom deny action for a specific configuration.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to create a custom deny action using a JSON definition
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_custom_deny" "custom_deny" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  custom_deny = file("${path.module}/custom_deny.json")
}

output "custom_deny_id" {
  value = akamai_appsec_custom_deny.custom_deny.custom_deny_id
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `custom_deny` - (Required) The JSON-formatted definition of the custom deny action ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#63df3de3)).

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

*`custom_deny_id` - The ID of the new custom deny action.

