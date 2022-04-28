---
layout: "akamai"
page_title: "Akamai: KRS Attack Groups"
subcategory: "Application Security"
description: |-
 KRS Attack Groups
---


# akamai_appsec_attack_groups

**Scopes**: Security policy; attack group

Returns the action and the condition-exception information for an attack group or set of attack groups. Attack groups are collections of Kona Rule Set rules used to streamline the management of website protections.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/attack-groups](https://techdocs.akamai.com/application-security/reference/get-policy-attack-groups)

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

// USE CASE: User wants to view the action and the condition-exception information for an attack group.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_attack_groups" "attack_group" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  attack_group       = "SQL"
}
output "attack_group_action" {
  value = data.akamai_appsec_attack_groups.attack_group.attack_group_action
}
output "condition_exception" {
  value = data.akamai_appsec_attack_groups.attack_group.condition_exception
}
output "json" {
  value = data.akamai_appsec_attack_groups.attack_group.json
}
output "output_text" {
  value = data.akamai_appsec_attack_groups.attack_group.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the attack group.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the attack group.
- `attack_group` (Optional). Unique name of the attack group you want to return information for. If not included, information is returned for all your attack groups.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `attack_group_action`. Action taken anytime the attack group is triggered. Valid values are:
  - **alert**. The event is recorded.
  - **deny**. The request is blocked.
  - **deny_custom_{custom_deny_id}**. The action defined by the custom deny is taken.
  - **none**. No action is taken.
- `condition_exception`. Conditions and exceptions assigned to the attack group.
- `json`. JSON-formatted list of the action and the condition-exception information for the attack group. This option is available only if the `attack_group` argument is included in the Terraform configuration file.
- `output_text`. Tabular report showing the attack group's action as well as Boolean values indicating whether conditions and exceptions have been configured for the group.