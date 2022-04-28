---
layout: "akamai"
page_title: "Akamai: Evaluation Attack Group"
subcategory: "Application Security"
description: |-
 Eval Group
---

# akamai_appsec_eval_group

**Scopes**: Evaluation attack group

Modifies the action and the conditions and exceptions for an evaluation mode attack group.

Note that this resource is only available to organizations running the Adaptive Security Engine (ASE) beta. For more information about ASE, please contact your Akamai representative.

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

// USE CASE: User wants to add an action and condition-exception information to an evaluation attack group by using a JSON input file.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_eval_group" "eval_attack_group" {
  config_id           = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id  = "gms1_134637"
  attack_group        = "SQL"
  attack_group_action = "deny"
  condition_exception = file("${path.module}/condition_exception.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration where evaluation is taking place.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the evaluation process.
- `attack_group` (Required). Unique identifier of the evaluation attack group being modified.
- `attack_group_action` (Required). Action to be taken any time the attack group is triggered. Allowed values are:
  - **alert**. Record the event.
  - **deny**. Block the request.
  - **deny_custom_{custom_deny_id}**. Take the action specified by the custom deny.
  - **none**. Take no action.
- `condition_exception` (Optional). Path to a JSON file containing properties and property values for the attack group. For more information, the [Modify the exceptions of an attack group](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putattackgroupconditionexception) section of the Application Security API documentation.