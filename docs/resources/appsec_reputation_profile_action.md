---
layout: "akamai"
page_title: "Akamai: Reputation Profile Action"
subcategory: "Application Security"
description: |-
 Reputation Profile Action
---

# akamai_appsec_reputation_profile_action

**Scopes**: Reputation profile

Modifies the action taken when a reputation profile is triggered.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/reputation-profiles/{reputationProfileId}](https://techdocs.akamai.com/application-security/reference/put-reputation-profile-action)

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

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_reputation_profile_action" "appsec_reputation_profile_action" {
  config_id             = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id    = "gms1_134637"
  reputation_profile_id = 130713
  action                = "alert"
}

output "reputation_profile_id" {
  value = akamai_appsec_reputation_profile_action.appsec_reputation_profile_action.reputation_profile_id
}

output "reputation_profile_action" {
  value = akamai_appsec_reputation_profile_action.appsec_reputation_profile_action.action
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the reputation profile action being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the reputation profile action being modified.
- `reputation_profile_id` (Required). Unique identifier of the reputation profile whose action is being modified.
- `action` (Required). Action taken any time the reputation profile is triggered. Allows values are:
  - **alert**. Record the event.
  - **deny**. Block the request.
  - **deny_custom_{custom_deny_id}**. Take the action specified by the custom deny.
  - **none**. Take no action.