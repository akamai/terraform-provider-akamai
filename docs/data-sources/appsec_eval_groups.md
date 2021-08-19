---
layout: "akamai"
page_title: "Akamai: Evaluation Attack Groups"
subcategory: "Application Security"
description: |-
 Evaluation Attack Groups
---

# akamai_appsec_eval_groups

Use the `akamai_appsec_eval_groups` data source to list the action and condition-exception information for an evaluation attack
group or groups.
__BETA__ This is Adaptive Security Engine(ASE) related data source. Please contact your akamai representative if you want to learn more

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
data "akamai_appsec_eval_groups" "eval_attack_group" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  attack_group = var.attack_group
}
output "eval_attack_group_action" {
  value = data.akamai_appsec_eval_groups.eval_attack_group.attack_group_action
}
output "condition_exception" {
  value = data.akamai_appsec_eval_groups.eval_attack_group.condition_exception
}
output "json" {
  value = data.akamai_appsec_eval_groups.eval_attack_group.json
}
output "output_text" {
  value = data.akamai_appsec_eval_groups.eval_attack_group.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `attack_group` - (Optional) The ID of the eval attack group to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `attack_group_action` - The eval attack group's action, either `alert`, `d

eny`, or `none`.

* `condition_exception` - The eval attack group's conditions and exceptions.

* `json` - A JSON-formatted list of the action and condition-exception information for the specified eval attack
group. This output is only generated if an attack group is specified.

* `output_text` - A tabular display showing, for the specified eval attack group or groups, the eval attack group's action and
boolean values indicating whether conditions and exceptions are present.
