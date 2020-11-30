---
layout: "akamai"
page_title: "Akamai: KRS Rule Actions"
subcategory: "Application Security"
description: |-
 KRS Rule Actions
---

# akamai_appsec_rule_actions

(Beta) Use the `akamai_appsec_rule_actions` data source to retrieve the action taken for each rule for a given security configuration version and security policy, or for a specific rule.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view all rule actions
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_rule_actions" "rule_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  policy_id = var.policy_id
}
output "rule_actions_text" {
  value = data.akamai_appsec_rule_actions.rule_actions.output_text
}
output "rule_actions_json" {
  value = data.akamai_appsec_rule_actions.rule_actions.json
}

// USE CASE: user wants to view a rule action
data "akamai_appsec_rule_actions" "rule_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  policy_id = var.policy_id
  rule_id = var.rule_id
}
output "rule_action_action" {
  value = akamai_appsec_rule_actions.rule_action.action
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `rule_id` - (Optional) The ID of a specific rule.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display of the ID and action for all rules in the security policy.

* `json` - A JSON-formatted display of the ID and action for all rules in the security policy.

* `action` - The action configured for the given rule if a `rule_id` was specified.
