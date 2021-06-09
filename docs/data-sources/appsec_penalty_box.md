---
layout: "akamai"
page_title: "Akamai: Penalty Box Settings"
subcategory: "Application Security"
description: |-
 Penalty Box
---

# akamai_appsec_penalty_box

Use the `akamai_appsec_penalty_box` data source to retrieve the penalty box settings for a specified security policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view penalty box settings
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_penalty_box" "penalty_box" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}

output "penalty_box_action" {
  value = data.akamai_appsec_penalty_box.penalty_box.action
}

output "penalty_box_enabled" {
  value = data.akamai_appsec_penalty_box.penalty_box.enabled
}

output "penalty_box_text" {
  value = data.akamai_appsec_penalty_box.penalty_box.output_text
}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `action` - The action for the penalty box: `alert`, `deny`, or `none`.

* `enabled` - Either `true` or `false`, indicating whether penalty box protection is enabled.

* `output_text` - A tabular display of the `action` and `enabled` information.

