---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_threat_intel

**Scopes**: Security policy

Enables or disables threat intelligence for a security policy. This resource is only available to organizations running the Adaptive Security Engine (ASE) beta Please contact your Akamai representative for more information.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/rules/threat-intel](https://techdocs.akamai.com/application-security/reference/put-rules-threat-intel)

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

// USE CASE: User wants to update the threat intelligence setting.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_threat_intel" "threat_intel" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  threat_intel       = "on"
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the threat intelligence protection settings being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the threat intelligence protection settings being modified.
- `threat_intel` (Required). Set to `on` to enable threat intelligence protection; set to **off** to disable threat intelligence protection.