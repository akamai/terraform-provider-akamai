---
layout: "akamai"
page_title: "Akamai: Eval Rule Action"
subcategory: "Application Security"
description: |-
 Eval Rule Action
---

# akamai_appsec_eval_rule_action

Use the `akamai_appsec_eval_rule_action` resource to update the action for a specific rule you want to evaluate.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// OPEN API --> https://developer.akamai.com/api/cloud_security/application_security/v1.html#putevalrule

// USE CASE: user wants to set the eval rule action
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_eval_rule_action" "rule_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  rule_id = var.rule_id
  rule_action = var.action
}

//TF destroy - set the action to none.

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `rule_id` - (Required) The ID of the rule being evaluated.

* `action` - (Required) The action to be taken: `alert` to record the trigger of the event, `deny` to block the request, or `none` to take no action.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None


