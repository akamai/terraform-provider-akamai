---
layout: "akamai"
page_title: "Akamai: Evaluation Attack Groups"
subcategory: "Application Security"
description: |-
 Evaluation Attack Groups
---



# akamai_appsec_eval_groups

**Scopes**: Security policy; evaluation attack group

Returns the action and the condition-exception information for an evaluation attack group or collection of groups. Note that this data source is only available to organizations running the Adaptive Security Engine (ASE) beta. For more information on ASE, please contact your Akamai representative.

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

// USE CASE: User wants to add an action and condition-exception information to an evaluation attack group using a JSON-formatted input file.

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

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the evaluation attack group.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the evaluation attack group.
- `attack_group` (Optional). Unique identifier of the evaluation attack group you want to return information for. If not included, information is returned for all your evaluation attack groups.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `attack_group_action`. Action taken anytime the attack group is triggered. Valid values are:
  - **alert**. Record the event.
  - **deny**. Block the request
  - **deny_custom_{custom_deny_id}**. The action defined by the custom deny is taken.
  - **none**. Take no action.
- `condition_exception`. Conditions and exceptions associated with the attack group.
- `json`. JSON-formatted list of the action and the condition-exception information for the specified attack group. This option is only available if the `attack_group` argument is included in the Terraform configuration.
- `output_text`. Tabular report showing the attack group's action as well as Boolean values indicating whether conditions and exceptions have been configured for the group.

