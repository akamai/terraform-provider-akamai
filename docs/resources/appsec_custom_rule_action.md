---
layout: "akamai"
page_title: "Akamai: CustomRuleAction"
subcategory: "Application Security"
description: |-
  Custom Rule Action
---

# akamai_appsec_custom_rule_action

**Scopes**: Custom rule

Associates an action with a custom rule. Custom rules are rules that you define yourself and are not part of the Kona Rule Set.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/custom-rules](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putactionruleid)

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

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_custom_rule_action" "create_custom_rule_action" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  custom_rule_id     = 12345
  custom_rule_action = "alert"
}

output "custom_rule_id" {
  value = akamai_appsec_custom_rule_action.create_custom_rule_action.custom_rule_id
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom rule action being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the custom rule action being modified d.
- `custom_rule_id` (Required). Unique identifier of the custom rule whose action is being modified.
- `custom_rule_action` (Required). Action to be taken when the custom rule is invoked. Allowed values are:
  - **alert**. Record the event.
  - **deny**. Block the request.
  - **deny_custom_{custom_deny_id}**. Take the action specified by the custom deny.
  - **none**. Take no action.

