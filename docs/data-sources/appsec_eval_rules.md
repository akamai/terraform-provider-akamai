---
layout: "akamai"
page_title: "Akamai: KRS Eval Rule
subcategory: "Application Security"
description: |-
 KRS Eval Rules
---

# akamai_appsec_eval_rules

Use the `akamai_appsec_eval_rules` data source to list the action and condition-exception information
for a rule or rules you want to evaluate.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view action and condition-exception information for an evaluation rule
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_eval_rules" "eval_rule" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  rule_id = var.rule_id
}
output "eval_rule_action" {
  value = akamai_appsec_eval_rules.eval_rule.eval_rule_action
}
output "condition_exception" {
  value = akamai_appsec_eval_rules.eval_rule.condition_exception
}
output "json" {
  value = akamai_appsec_eval_rules.eval_rule.json
}
output "output_text" {
  value = akamai_appsec_eval_rules.eval_rule.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `rule_id` - (Optional) The ID of the rule to use. If not specified, information about all rules will be returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `eval_rule_action` - The eval rule's action, either `alert`, `deny`, or `none`.

* `condition_exception` - The eval rule's conditions and exceptions.

* `json` - A JSON-formatted list of the action and condition-exception information for the specified eval rule.
This output is only generated if an eval rule is specified.

* `output_text` - A tabular display showing, for the specified eval rule or rules, the rule action and boolean values
indicating whether conditions and exceptions are present.

