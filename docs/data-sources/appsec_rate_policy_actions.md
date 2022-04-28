---
layout: "akamai"
page_title: "Akamai: Rate Policy Actions"
subcategory: "Application Security"
description: |-
 Rate Policy Actions
---

# akamai_appsec_rate_policy_actions

**Scopes**: Security policy; rate policy

Returns information about your rate policy actions. Actions specify what happens any time a rate policy is triggered: the issue could be ignored, the request could be denied, or an alert could be generated.

**Related API Endpoint:** [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/rate-policies](https://techdocs.akamai.com/application-security/reference/get-rate-policies-actions)

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

// USE CASE: User wants to view all the rate policy actions associated with a security policy.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_rate_policy_actions" "rate_policy_actions" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}
output "rate_policy_actions" {
  value = data.akamai_appsec_rate_policy_actions.rate_policy_actions.output_text
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the rate policies and rate policy actions.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the rate policies and rate policy actions.
- `rate_policy_id` (Optional). Unique identifier of the rate policy you want to return action information for. If not included, action information is returned for all your rate policies.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the ID, IPv4 action, and IPv6 action of the rate policies.