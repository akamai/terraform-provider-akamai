---
layout: "akamai"
page_title: "Akamai: Attack Group"
subcategory: "Application Security"
description: |-
 Attack Group
---

# akamai_appsec_attack_group

**Scopes**: Attack group

Modify an attack group's action, conditions, and exceptions. Attack groups are collections of Kona Rule Set rules used to streamline the management of website protections.

**Related API Endpoints**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/attack-groups/{attackGroupId}](https://techdocs.akamai.com/application-security/reference/put-attack-group-condition-exception) *and* [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/attack-groups/{attackGroupId}/condition-exception](https://techdocs.akamai.com/application-security/reference/put-attack-group-condition-exception)

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

// USE CASE: User wants to add action and condition-exception information to an attack group by using a JSON input file.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_attack_group" "attack_group" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id  = "gms1_134637"
  attack_group        = "SQL"
  attack_group_action = "deny"
  condition_exception = file("${path.module}/condition_exception.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the attack group being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the attack group being modified.
- `attack_group` (Required). Unique name of the attack group being modified.
- `attack_group_action` (Required). Action taken any time the attack group is triggered. Allowed values are:
  - **alert**. Record information about the request.
  - **deny**. Block the request,
  - **deny_custom_{custom_deny_id}**. Take the action specified by the custom deny.
  - **none**. Take no action.
- `condition_exception` (Optional). Path to a JSON file containing the conditions and exceptions to be assigned to the attack group. You can view a sample JSON file in the [Modify the exceptions of an attack group](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putattackgroupconditionexception) section of the Application Security API documentation.