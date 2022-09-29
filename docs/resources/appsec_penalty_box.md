---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_penalty_box

**Scopes**: Security policy

Modifies the penalty box settings for a security policy.
When the penalty box is enabled for a policy, clients that trigger a WAF Deny action are placed in the “penalty box”.
There, the action you select for penalty box (either Alert or Deny ) continues to apply to any requests from that client for the next 10 minutes.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/match-targets/sequence](https://techdocs.akamai.com/application-security/reference/put-policy-penalty-box)

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
resource "akamai_appsec_penalty_box" "penalty_box" {
  config_id              = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id     = "gms1_134637"
  penalty_box_protection = true
  penalty_box_action     = "deny"
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the penalty box settings being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the penalty box settings being modified.
- `penalty_box_protection` (Required). Set to **true** to enable penalty box protection; set to **false** to disable penalty box protection.
- `penalty_box_action` (Required). Action taken any time penalty box protection is triggered. Allowed values are:
  - **alert**. Record the event.
  - **deny**. Block the request.
  - **deny_custom_{custom_deny_id}**. Take the action specified by the custom deny.
  - **none**. Take no action.
