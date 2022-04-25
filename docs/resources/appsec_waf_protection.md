---
layout: "akamai"
page_title: "Akamai: WAF Protection"
subcategory: "Application Security"
description: |-
 WAF Protection
---

# akamai_appsec_waf_protection

**Scopes**: Security policy

Enables or disables Web Application Firewall (WAF) protection for a security policy.

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

// USE CASE: User wants to enable or disable WAF protection.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_waf_protection" "protection" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  enabled            = true
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the WAF protection settings being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the WAF protection settings being modified.
- `enabled` (Required). Set to **true** to enable WAF protection; set to **false** to disable WAF protection.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the current protection settings.