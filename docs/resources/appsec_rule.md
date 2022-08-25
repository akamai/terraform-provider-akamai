---
layout: "akamai"
page_title: "Akamai: Rule"
subcategory: "Application Security"
description: |-
 Rule
---

# akamai_appsec_rule

**Scopes**: Rule

Modifies a Kona Rule Set rule's action, conditions, and exceptions.

**Related API Endpoints**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/rules/{ruleId}](https://techdocs.akamai.com/application-security/reference/put-rule) *and* [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/rules/{ruleId}/condition-exception](https://techdocs.akamai.com/application-security/reference/put-rule-condition-exception)

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

// USE CASE: User wants to add an action and condition-exception information to a rule by using a JSON-formatted input file.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_rule" "rule" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id  = "gms1_134637"
  rule_id             = 60029316
  rule_action         = "deny"
  condition_exception = file("${path.module}/condition_exception.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the Kona Rule Set rule being modified.

- `security_policy_id` (Required). Unique identifier of the security policy associated with the Kona Rule Set rule being modified.

- `rule_id` (Required). Unique identifier of the rule being modified.

- `rule_action` - (Required except when the policy in ASE AUTO mode) Allowed values are:
  - **alert**. Record the event.
  - **deny**. Block the request.
  - **deny_custom_{custom_deny_id}**. Take the action specified by the custom deny.
  - **none**. Take no action. or `none` to take no action.

 __ASE Beta__. if policy is in `ASE_AUTO` mode, only condition_exception can be modified, "ASE" (Adaptive Security Engine) is currently in beta. Please contact your Akamai representative to learn more.

- `condition_exception` (Optional). Path to a JSON file containing a description of the conditions and exceptions to be associated with a rule. 