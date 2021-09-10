---
layout: "akamai"
page_title: "Akamai: Reputation Profile Actions"
subcategory: "Application Security"
description: |-
 Reputation Profile Actions
---

## akamai_appsec_reputation_profile_actions

**Scopes**: Security policy; reputation profile

Returns action information for your reputation profiles. Actions specify what happens any time a profile is triggered: the issue could be ignored, the request could be denied, or an alert could be generated.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/reputation-profiles](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getreputationprofileactions)

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

// USE CASE: User wants to view the reputation profile actions associated with a security policy.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_reputation_profile_actions" "reputation_profile_actions" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}
output "reputation_profile_actions_text" {
  value = data.akamai_appsec_reputation_profile_actions.reputation_profile_actions.output_text
}
output "reputation_profile_actions_json" {
  value = data.akamai_appsec_reputation_profile_actions.reputation_profile_actions.json
}

// USE CASE: User wants to view the action for a specific reputation profile.

data "akamai_appsec_reputation_profile_actions" "reputation_profile_actions2" {
  config_id             = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id    = "gms1_134637"
  reputation_profile_id = "12345"
}

output "reputation_profile_actions2" {
  value = data.akamai_appsec_reputation_profile_actions.reputation_profile_actions.action
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the reputation profiles.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the reputation profiles.
- `reputation_profile_id` (Optional). Unique identifier of the reputation profile you want to return information for. If not included, information is returned for all your reputation profiles.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `action`. Action taken any time the reputation profile is triggered. Valid values are:
  - **alert**. Record the event.
  - **deny**. Block the request.
  - **deny_custom_{custom_deny_id}**. The action defined by the custom deny is taken.
  - **none**. Take no action.
- `json`. JSON-formatted report of the reputation profile action information.
- `output_text`. Tabular report of the reputation profile action information.

