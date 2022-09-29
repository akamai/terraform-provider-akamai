---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_rules

**Scopes**: Security policy; rule

Returns the action and the condition-exception information for your Kona Rule Set (KRS) rules.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/rules](https://techdocs.akamai.com/application-security/reference/get-policy-rules)

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

// USE CASE: User wants to view the action and the condition-exception information for a rule.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_rules" "rule" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  rule_id            = "60029316"
}
output "rule_action" {
  value = data.akamai_appsec_rules.rule.rule_action
}
output "condition_exception" {
  value = data.akamai_appsec_rules.rule.condition_exception
}
output "json" {
  value = data.akamai_appsec_rules.rule.json
}
output "output_text" {
  value = data.akamai_appsec_rules.rule.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the rules.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the rules.
- `rule_id` (Optional). Unique identifier of the Kona Rule Set rule you want to return information for. If not included, information is returned for all your KRS rules.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `rule_action`. Action taken anytime the rule is triggered. Valid values are:
  - **alert**. The event is recorded.
  - **deny**. The request is blocked.
  - **deny_custom_{custom_deny_id}**. The action defined by the custom deny is taken.
  - **none**. No action is taken.
- `condition_exception`. Conditions and exceptions associated with the rule.
- `json`. JSON-formatted list of the action and the condition-exception information for the rule. This option is only available if the `rule_id` argument is included in your Terraform configuration file.
- `output_text`. Tabular report showing the rule action as well as Boolean values indicating whether conditions and exceptions are configured.