---
layout: "akamai"
page_title: "Akamai: KRS Rules"
subcategory: "Application Security"
description: |-
 KRS Rules
---

# akamai_appsec_rules

Use the `akamai_appsec_rules` data source to list the action and condition-exception information for a rule or rules.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view action and condition-exception information for a rule
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_rules" "rule" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  rule_id = var.rule_id
}
output "rule_action" {
  value = akamai_appsec_rules.rule.rule_action
}
output "condition_exception" {
  value = akamai_appsec_rules.rule.condition_exception
}
output "json" {
  value = akamai_appsec_rules.rule.json
}
output "output_text" {
  value = akamai_appsec_rules.rule.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `rule_id` - (Optional) The ID of the rule to use. If not specified, information about all rules will be returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `rule_action` - The rule's action, either `alert`, `deny`, or `none`.

* `condition_exception` - The rule's conditions and exceptions.

* `json` - A JSON-formatted list of the action and condition-exception information for the specified rule.
This output is only generated if a rule is specified.

* `output_text` - A tabular display showing, for the specified rule or rules, the rule action and boolean values
indicating whether conditions and exceptions are present.

