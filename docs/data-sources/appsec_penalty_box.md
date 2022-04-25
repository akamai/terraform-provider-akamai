---
layout: "akamai"
page_title: "Akamai: Penalty Box Settings"
subcategory: "Application Security"
description: |-
 Penalty Box
---


# akamai_appsec_penalty_box

**Scopes**: Security policy

Returns penalty box settings for the specified security policy. When using automated attack groups, and when the penalty box is enabled, clients that trigger an attack group are placed in the “penalty box.” That means that, for the next 10 minutes, all requests from that client are ignored.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/penalty-box](https://techdocs.akamai.com/application-security/reference/get-policy-penalty-box)

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
data "akamai_appsec_penalty_box" "penalty_box" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "penalty_box_action" {
  value = data.akamai_appsec_penalty_box.penalty_box.action
}

output "penalty_box_enabled" {
  value = data.akamai_appsec_penalty_box.penalty_box.enabled
}

output "penalty_box_text" {
  value = data.akamai_appsec_penalty_box.penalty_box.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the penalty box settings.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the penalty box settings.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `action`. Action taken any time the penalty box is triggered. Valid values are:
  - **alert**. Record the event.
  - **deny**. The request is blocked.
  - **deny_custom_{custom_deny_id}**. The action defined by the custom deny is taken.
  - **none**. Take no action.
- `enabled`. If **true**, penalty box protection is enabled. If **false**, penalty box protection is disabled.
- `output_text`. Tabular report of penalty box protection settings.