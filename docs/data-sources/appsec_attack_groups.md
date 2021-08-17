---
layout: "akamai"
page_title: "Akamai: KRS Attack Groups"
subcategory: "Application Security"
description: |-
 KRS Attack Groups
---

# akamai_appsec_attack_groups

Use the `akamai_appsec_attack_groups` data source to list the action and condition-exception information for an attack
group or groups.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view action and condition-exception information for an attack group
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_attack_groups" "attack_group" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  attack_group = var.attack_group
}
output "attack_group_action" {
  value = akamai_appsec_attack_groups.attack_group.attack_group_action
}
output "condition_exception" {
  value = akamai_appsec_attack_groups.attack_group.condition_exception
}
output "json" {
  value = akamai_appsec_attack_groups.attack_group.json
}
output "output_text" {
  value = akamai_appsec_attack_groups.attack_group.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `attack_group` - (Optional) The ID of the attack group to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `attack_group_action` - The attack group's action, either `alert`, `deny`, or `none`.

* `condition_exception` - The attack group's conditions and exceptions.

* `json` - A JSON-formatted list of the action and condition-exception information for the specified attack
group. This output is only generated if an attack group is specified.

* `output_text` - A tabular display showing, for the specified attack group or groups, the attack group's action and
boolean values indicating whether conditions and exceptions are present.

