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

**Related API Endpoints**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/rules/{ruleId}](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putruleaction) *and* [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/rules/{ruleId}/condition-exception](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putruleconditionexception)

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

- `rule_action` (Optional). Action to be taken anytime the rule is triggered. Allowed values are:

  - **alert**. Record the event.
  - **deny**. Block the request.
  - **deny_custom_{custom_deny_id}** Take the action specified by the custom deny.
  - **none**. Take no action.

  If you are running the Adaptive Security Engine (ASE) beta in **ASE_AUTO** mode, you can't modify the rule action.
  You can only modify the rule's conditions and exception. 
  Please contact your Akamai representative for more information.

- `condition_exception` (Optional). Path to a JSON file containing a description of the conditions and exceptions to be associated with a rule. You can view a sample JSON file in the [Modify the conditions and exceptions of a rule](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putruleconditionexception) section of the Application Security API documentation.

