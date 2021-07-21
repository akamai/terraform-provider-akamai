---
layout: "akamai"
page_title: "Akamai: Attack Group"
subcategory: "Application Security"
description: |-
 Attack Group
---

# akamai_appsec_attack_group

Use the `akamai_appsec_attack_group` resource to create or modify an attack group's action, conditions and exceptions. When the conditions are met, the rule’s actions are ignored and not applied to that specific traffic.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to add action and condition-exception information to an attack group using a JSON input file
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_attack_group" "attack_group" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  attack_group = var.attack_group
  attack_group_action = var.attack_group_action
  condition_exception = file("${path.module}/condition_exception.json")
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `attack_group` - The attack group to use.

* `attack_group_action` - (Required) The action to be taken: `alert` to record the trigger of the event, `deny` to block the request, `deny_custom_{custom_deny_id}` to execute a custom deny action, or `none` to take no action.

* `condition_exception` - (Optional) The name of a file containing a JSON-formatted description of the conditions and exceptions to use ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putattackgroupconditionexception)).

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

