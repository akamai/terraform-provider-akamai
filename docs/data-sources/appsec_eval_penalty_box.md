---
layout: "akamai"
page_title: "Akamai: Penalty Box Settings"
subcategory: "Application Security"
description: |-
 Penalty Box
---


# akamai_appsec_eval_penalty_box

**Scopes**: Security policy

 __ASE_Beta__.:
Returns the penalty box settings for a security policy in evaluation mode - evaluation penalty box. 
When the penalty box is enabled for a policy in evaluation mode, clients that trigger a WAF Deny action are placed in the “penalty box”.
There, the action you select for the penalty box (either Alert or Deny) continues to apply to any requests from that client for the next 10 minutes.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/eval_penalty-box](https://techdocs.akamai.com/application-security/reference/get-policy-eval_penalty-box)

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

// USE CASE: User wants to view penalty box settings.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_eval_penalty_box" "eval_penalty_box" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "eval_penalty_box_action" {
  value = data.akamai_appsec_eval_penalty_box.eval_penalty_box.action
}

output "eval_penalty_box_enabled" {
  value = data.akamai_appsec_eval_penalty_box.eval_penalty_box.enabled
}

output "eval_penalty_box_text" {
  value = data.akamai_appsec_eval_penalty_box.eval_penalty_box.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the evaluation penalty box settings.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the evaluation penalty box settings.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `action`. Action taken any time the penalty box is triggered. Valid values are:
  - **alert**. Record the event.
  - **deny**. The request is blocked.
  - **deny_custom_{custom_deny_id}**. The action defined by the custom deny is taken.
  - **none**. Take no action.
- `enabled`. If **true**, evaluation penalty box protection is enabled. If **false**, evaluation penalty box protection is disabled.
- `output_text`. Tabular report of evaluation penalty box protection settings.
