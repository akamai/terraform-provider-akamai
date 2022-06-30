---
layout: "akamai"
page_title: "Akamai: Rate Policy Action"
subcategory: "Application Security"
description: |-
  Rate Policy Action
---

# akamai_appsec_rate_policy_action

**Scopes**: Rate policy

Creates, modifies, or deletes the actions associated with a rate policy.
By default, rate policies take no action when triggered.
Note that you must set separate actions for requests originating from an IPv4 IP address and for requests originating from an IPv6 address.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/rate-policies/{ratePolicyId}](https://techdocs.akamai.com/application-security/reference/put-rate-policy-action)

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
// USE CASE: User wants to create a rate policy and rate policy actions for a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_rate_policy" "appsec_rate_policy" {
  config_id   = data.akamai_appsec_configuration.configuration.config_id
  rate_policy = file("${path.module}/rate_policy.json")
}
resource "akamai_appsec_rate_policy_action" "appsec_rate_policy_action" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  rate_policy_id     = akamai_appsec_rate_policy.appsec_rate_policy.rate_policy_id
  ipv4_action        = "deny"
  ipv6_action        = "deny"
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the rate policy action being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the rate policy whose action is being modified.
- `rate_policy_id` (Required). Unique identifier of the rate policy whose action is being modified.
- `ipv4_action` (Required). Rate policy action for requests coming from an IPv4 IP address. Allowed actions are:
  - **alert**. Record the event.
  - **deny**. Block the request.
  - **deny_custom{custom_deny_id}**. Take the action specified by the custom deny.
  - **none**. Take no action.
- `ipv6_action` (Required). Rate policy action for requests coming from an IPv6 IP address. Allowed actions are:
  - **alert**. Record the event.
  - **deny**. Block the request.
  - **deny_custom{custom_deny_id}**. Take the action specified by the custom deny.
