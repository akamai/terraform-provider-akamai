---
layout: "akamai"
page_title: "Akamai: Slowpost Protection"
subcategory: "Application Security"
description: |-
 Slowpost Protection
---

# akamai_appsec_slowpost_protection

**Scopes**: Security policy

Enables or disables slow POST protection for a security configuration and security policy. Slow POST protections help defend a site against attacks that try to tie up the site by using extremely slow requests and responses.

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

// USE CASE: User wants to enable or disable slow post protections.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_slowpost_protection" "protection" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  enabled            = true
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the slow POST protection settings being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the slow POST protection settings being modified.
- `enabled` (Required). Set to **true** to enable slow POST protection; set to **false** to disable slow POST protection.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the current protection settings.