---
layout: "akamai"
page_title: "Akamai: Eval Rule"
subcategory: "Application Security"
description: |-
 Eval Rule
---

# akamai_appsec_eval_rule

**Scopes**: Evaluation rule

Creates or modifies an evaluation rule's action, conditions, and exceptions.
Evaluation rules are Kona Rule Set rules used when running a security configuration in evaluation mode.
Changes to these rules do not affect the rules used on your production network.

**Related API Endpoints**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/eval-rules/{ruleId}](https://techdocs.akamai.com/application-security/reference/put-policy-eval-rule) *and* [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/eval-rules/{ruleId}/condition-exception](https://techdocs.akamai.com/application-security/reference/put-condition-exception)

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

// USE CASE: User wants to add an action and condition-exception information to an evaluation rule by using a JSON input file.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_eval_rule" "eval_rule" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id  = "gms1_134637"
  rule_id             = 60029316
  rule_action         = "deny"
  condition_exception = file("${path.module}/condition_exception.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration in evaluation mode.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the evaluation process.
- `rule_id` (Required). Unique identifier of the evaluation rule being modified.
- `rule_action` (Required). Action to be taken any time the evaluation rule is triggered, Allowed actions are:
  - **alert**. Record the event.
  - **deny**. Block the request.
  - **deny_custom_{custom_deny_id}**. Take the action specified by the custom deny.
  - **none**. Take no action.
- `condition_exception` (Optional). Path to a JSON file containing the conditions and exceptions to be applied to the evaluation rule. To view a sample JSON file, see the [Modify the conditions and exceptions for an evaluation rule](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevalconditionsexceptions) section of the Application Security API documentation.