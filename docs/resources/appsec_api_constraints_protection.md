---
layout: "akamai"
page_title: "Akamai: API Constraints Protection"
subcategory: "Application Security"
description: |-
 API Constraints Protection
---

# akamai_appsec_api_constraints_protection

Use the `akamai_appsec_api_constraints_protection` resource to enable or disable API constraints protection for a given configuration and security policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to enable or disable API constraints protection
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

resource "akamai_appsec_api_constraints_protection" "protection" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  enabled = var.enabled
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `enabled` - (Required) Whether to enable API constraints protection: either `true` or `false`.


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the current protection settings.

