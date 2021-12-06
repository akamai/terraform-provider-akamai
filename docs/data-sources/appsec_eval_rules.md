---
layout: "akamai"
page_title: "Akamai: Evaluation Rule"
subcategory: "Application Security"
description: |-
 Evaluation Rules
---


# akamai_appsec_eval_rules

**Scopes**: Security policy; evaluation rule

Returns the action and the condition-exception information for a rule or set of rules being used in evaluation mode.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/eval-rules](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getevalrules)

## Example Usage

Basic usage:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: User wants to view the action and the condition-exception information for an evaluation rule.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_eval_rules" "eval_rule" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  rule_id            = "60029316"
}
output "eval_rule_action" {
  value = data.akamai_appsec_eval_rules.eval_rule.eval_rule_action
}
output "condition_exception" {
  value = data.akamai_appsec_eval_rules.eval_rule.condition_exception
}
output "json" {
  value = data.akamai_appsec_eval_rules.eval_rule.json
}
output "output_text" {
  value = data.akamai_appsec_eval_rules.eval_rule.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration running in evaluation mode.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the evaluation rule.
- `rule_id` (Optional). Unique identifier of the evaluation rule you want to return information for. If not included, information is returned for all your evaluation rules.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `eval_rule_action`. Action taken anytime the evaluation rule is triggered. Valid values are:
  - **alert**. Record the event,
  - **deny**. Reject the request.
  - **deny_custom_{custom_deny_id}**. The action defined by the custom deny is taken.
  - **none**. Take no action.
- `condition_exception`. Conditions and exceptions associated with the rule.
- `json`. JSON-formatted list of the action and the condition-exception information for the rule. This output is only generated if the `rule_id` argument is included.
- `output_text`. Tabular report showing the rule action as well as Boolean values indicating whether conditions and exceptions have been configured for the rule.

