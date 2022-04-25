---
layout: "akamai"
page_title: "Akamai: Reputation Protection"
subcategory: "Application Security"
description: |-
 Reputation Protection
---

# akamai_appsec_reputation_protection

**Scopes**: Security policy

Enables or disables reputation protection for a security configuration and security policy.
Reputation profiles grade the security risk of an IP address based on previous activities associated with that address.
Depending on the reputation score and how your configuration has been set up, requests from a specific IP address can trigger an alert or even be blocked.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/protections](https://techdocs.akamai.com/application-security/reference/put-policy-protections)

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

// USE CASE: User wants to enable or disable reputation protections.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_reputation_protection" "protection" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  enabled            = true
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the reputation protection settings being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the reputation protection settings being modified.
- `enabled` (Required). Set to **true** to enable reputation protection; set to **false** to disable reputation protection.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the current protection settings.