---
layout: "akamai"
page_title: "Akamai: Penalty Box"
subcategory: "Application Security"
description: |-
 Penalty Box
---

# akamai_appsec_eval_penalty_box

**Scopes**: Security policy

 __ASE_Beta__.:
Modifies the penalty box settings for a security policy in evaluation mode - evaluation penalty box. 
When the penalty box is enabled for a policy in evaluation mode, clients that trigger a WAF Deny action are placed in the “penalty box”.
There, the action you select for the penalty box (either Alert or Deny) continues to apply to any requests from that client for the next 10 minutes.

**Related API Endpoint**:  [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/eval_penalty-box](https://techdocs.akamai.com/application-security/reference/put-policy-eval_penalty-box)
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

// USE CASE: User wants to update penalty box settings.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_eval_penalty_box" "eval_penalty_box" {
  config_id              = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id     = "gms1_134637"
  penalty_box_protection = true
  penalty_box_action     = "deny"
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the evaluation penalty box settings being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the evaluation penalty box settings being modified.
- `penalty_box_protection` (Required). Set to **true** to enable evaluation penalty box protection; set to **false** to disable evaluation penalty box protection.
- `penalty_box_action` (Required). Action taken any time evaluation penalty box protection is triggered. Allowed values are:
  - **alert**. Record the event.
  - **deny**. Block the request.
  - **deny_custom_{custom_deny_id}**. Take the action specified by the custom deny.
  - **none**. Take no action.
