---
layout: "akamai"
page_title: "Akamai: Eval Rule"
subcategory: "Application Security"
description: |-
 Eval Rule
---

# akamai_appsec_eval_rule

Use the `akamai_appsec_eval_rule` resource to create or modify an eval rule's action, conditions and exceptions. When the conditions are met, the rule’s actions are ignored and not applied to that specific traffic.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to add action and condition-exception information to an eval rule using a JSON input file
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_eval_rule" "eval_rule" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  rule_id = var.rule_id
  rule_action = var.action
  condition_exception = file("${path.module}/condition_exception.json")
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `rule_id` - (Required) The ID of the eval rule to use.

* `rule_action` - (Required) The action to be taken: `alert` to record the trigger of the event, `deny` to block the request, `deny_custom_{custom_deny_id}` to execute a custom deny action, or `none` to take no action.

* `condition_exception` - (Required) The name of a file containing a JSON-formatted description of the conditions and exceptions to use ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevalconditionsexceptions))

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

