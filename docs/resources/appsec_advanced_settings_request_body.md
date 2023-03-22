---
layout: akamai
subcategory: Application Security
---

# resource_akamai_appsec_advanced_settings_request_body

**Scopes**: Security configuration; security policy

The `resource_akamai_appsec_advanced_settings_request_body` resource allows you to update the Request Size Inspection Limit settings for a configuration.
This operation applies at the configuration level, and therefore applies to all policies within a configuration.
You may override this setting for a particular policy by specifying the policy using the security_policy_id parameter.

**Related API Endpoints**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/advanced-settings/request-body](https://techdocs.akamai.com/application-security/reference/put-policies-request-body)

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to set the Request Size Inspection Limit setting
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

resource "akamai_appsec_advanced_settings_request_body" "request_body" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  request_body_inspection_limit = 16
}

// USE CASE: user wants to override the Request Size Inspection Limit setting for a security policy
resource "akamai_appsec_advanced_settings_request_body" "policy_override" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  request_body_inspection_limit = 16
}
```

## Argument Reference

The following arguments are supported:

- `config_id` - (Required) The ID of the security configuration to use.

- `security_policy_id` - (Optional) The ID of a specific security policy to which the Request Size Inspection Limit setting should be applied. If not supplied, the indicated setting will be applied to all policies within the configuration.

- `request_body_inspection_limit` - (Required) Inspect request bodies up to a certain size. In exceptional cases, you can change the default value to a set limit: 'default', 8, 16, or 32 KB.
