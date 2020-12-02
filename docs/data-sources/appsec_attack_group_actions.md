---
layout: "akamai"
page_title: "Akamai: Attack Group Actions"
subcategory: "Application Security"
description: |-
 Attack Group Actions
---

# akamai_appsec_attack_group_actions

Use the `akamai_appsec_attack_group_actions` data source to retrieve a list of attack groups with their associated actions, or the action for a specific attack group.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view all attack group actions
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_attack_group_actions" "attack_group_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
}

output "attack_group_actions_text" {
  value = data.akamai_appsec_attack_group_actions.attack_group_actions.output_text
}

output "attack_group_actions_json" {
  value = data.akamai_appsec_attack_group_actions.attack_group_actions.json
}

// USE CASE: user wants to view an attack group  action
data "akamai_appsec_attack_group_actions" "attack_group_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  attack_group = var.attack_group
}

output "attack_group_action" {
  value = data.akamai_appsec_attack_group_actions.attack_group_action.action
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `attack_group` - (Optional) The attack group to use. If not supplied, information about all attack groups will be returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `action` - The attack group action for the attack group if one was specified: `alert`, `deny`, or `none`. If the action is none, the attack group is inactive in the security policy.

* `output_text` - A tabular display showing the `action` and `group` name for each attack group.

* `json` - The attack group information in JSON format.

