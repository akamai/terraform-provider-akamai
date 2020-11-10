---
layout: "akamai"
page_title: "Akamai: CustomRuleAction"
subcategory: "APPSEC"
description: |-
  CustomRuleAction
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
  version = data.akamai_appsec_configuration.configuration.latest_version
  policy_id = "crAP_75829"
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

* `version` - (Required) The version number of the security configuration to use.

* `custom_rule_action` - (Required) The action to be taken when the custom rule is invoked. Must be one of the following:
  * alert
  * deny
  * none

* `policy_id` - (Required) The 

* `custom_rule_id` - (Required)

## Attribute Reference

In addition to the arguments above, the following attribute is exported:

* `custom_rule_id` - The ID of the custom rule.

