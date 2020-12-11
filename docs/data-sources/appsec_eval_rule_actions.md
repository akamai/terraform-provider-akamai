---
layout: "akamai"
page_title: "Akamai: Eval Rule Actions"
subcategory: "Application Security"
description: |-
 Eval Rule Actions
---

# akamai_appsec_eval_rule_actions

Use the `akamai_appsec_eval_rule_actions` data source to retrieve the rules available for evaluation and their actions, or the action for a specific rule available for evaluation.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view all eval rule actions
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_eval_rule_actions" "rule_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
}
output "rule_actions_text" {
  value = data.akamai_appsec_eval_rule_actions.rule_actions.output_text
}
output "rule_actions_json" {
  value = data.akamai_appsec_eval_rule_actions.rule_actions.json
}

// USE CASE: user wants to view an eval rule action
data "akamai_appsec_eval_rule_actions" "rule_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  rule_id = var.rule_id
}
output "rule_action" {
  value = akamai_appsec_eval_rule_actions.rule_action.action
}
output "rule_id" {
  value = akamai_appsec_eval_rule_actions.rule_action.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `rule_id` - (Optional) The ID of a specific rule. If not supplied, information about all eval rules will be returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display of the ID and action for all rules in the security policy.

* `json` - A JSON-formatted display of the ID and action for all rules in the security policy.

* `action` - The action configured for the given rule if a `rule_id` was specified.
