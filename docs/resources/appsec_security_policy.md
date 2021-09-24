---
layout: "akamai"
page_title: "Akamai: SecurityPolicy"
subcategory: "Application Security"
description: |-
 SecurityPolicy
---

# akamai_appsec_security_policy

**Scopes**: Security configuration

Creates a new security policy. The resource enables you to:

- Create a new, “blank” security policy.
- Create a new policy preconfigured with the default security policy settings.
- Clone an existing security policy.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies](https://developer.akamai.com/api/cloud_security/application_security/v1.html#postsecuritypolicies)

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

// USE CASE: User wants to create a new security policy.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_security_policy" "security_policy_create" {
  config_id              = data.akamai_appsec_configuration.configuration.config_id
  default_settings       = true
  security_policy_name   = "Documentation Policy"
  security_policy_prefix = "gms1"
}

output "security_policy_create" {
  value = akamai_appsec_security_policy.security_policy_create.security_policy_id
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration to be associated with the new security policy.
- `security_policy_name` (Required). Name of the new security policy.
- `security_policy_prefix` (Required). Four-character alphanumeric string prefix used in creating the security policy ID.
- `default_settings` (Optional). Set to **true** to assign default setting values to the new policy; set to **false** to create a “blank” security policy. If not included, the new policy will be created using the default settings.
- `create_from_security_policy_id` (Optional). Unique identifier of the existing security policy that the new policy will be cloned from.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `security_policy_id`. ID of the newly-created security policy.

