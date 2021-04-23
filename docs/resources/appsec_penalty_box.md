---
layout: "akamai"
page_title: "Akamai: Penalty Box"
subcategory: "Application Security"
description: |-
 Penalty Box
---

# akamai_appsec_penalty_box

Use the `akamai_appsec_penalty_box` resource to update the penalty box settings for a given security policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to update the penalty box settings
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_penalty_box" "penalty_box" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  penalty_box_protection = true
  penalty_box_action = var.action
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `penalty_box_protection` - (Required) A boolean value indicating whether to enable penalty box protection.

* `penalty_box_action` - (Required) The action to take when penalty box protection is triggered: `alert` to record the trigger event, `deny` to block the request, `deny_custom_{custom_deny_id}` to execute a custom deny action, or `none` to take no action. Ignored if `penalty_box_protection` is set to `false`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

