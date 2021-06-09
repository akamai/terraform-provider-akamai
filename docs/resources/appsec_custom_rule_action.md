---
layout: "akamai"
page_title: "Akamai: CustomRuleAction"
subcategory: "Application Security"
description: |-
  Custom Rule Action
---

# akamai_appsec_custom_rule_action


The `akamai_appsec_custom_rule_action` resource allows you to associate an action to a custom rule.


## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

resource "akamai_appsec_custom_rule_action" "create_custom_rule_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "crAP_75829"
  custom_rule_id = 12345
  custom_rule_action = "alert"
}

output "custom_rule_id" {
  value = akamai_appsec_custom_rule_action.create_custom_rule_action.custom_rule_id
}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The security policy to use.

* `custom_rule_id` - (Required) The custom rule for which to apply the action.

* `custom_rule_action` - (Required) The action to take when the custom rule is invoked: `alert` to record the trigger event, `deny` to block the request, `deny_custom_{custom_deny_id}` to execute a custom deny action, or `none` to take no action.

## Attribute Reference

In addition to the arguments above, the following attribute is exported:

* None

