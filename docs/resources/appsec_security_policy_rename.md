---
layout: "akamai"
page_title: "Akamai: SecurityPolicyRename"
subcategory: "Application Security"
description: |-
  SecurityPolicyRename
---

# akamai_appsec_security_policy_rename

**Scopes**: Security policy

Renames an existing security policy. Note that you can only change the name of the policy: once issued, the security policy ID can't be modified.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}](https://techdocs.akamai.com/application-security/reference/put-policy)

## Example Usage

Basic usage:

```terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: User wants to rename a security policy.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_security_policy_rename" "security_policy_rename" {
  config_id            = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id   = "gms1_134637"
  security_policy_name = "Documentation and Training Policy"
}

```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the security policy being renamed.
- `security_policy_id` (Required). Unique identifier of the security policy being renamed.
- `security_policy_name` (Required). New name to be given to the security policy.