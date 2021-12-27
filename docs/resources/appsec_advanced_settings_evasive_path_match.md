---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsEvasivePathMatch"
subcategory: "Application Security"
description: |-
  AdvancedSettingsEvasivePathMatch
---

# resource_akamai_appsec_advanced_settings_evasive_path_match

The `resource_akamai_appsec_advanced_settings_evasive_path_match` resource allows you to enable, disable, or update the evasive path match setting for a configuration. This setting determines whether fuzzy matching is used to make URL matching more inclusive.
This operation applies at the configuration level, and therefore applies to all policies within a configuration. You may override this setting for a particular policy by specifying the policy using the security_policy_id parameter.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to set the evasive path match setting
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

resource "akamai_appsec_advanced_settings_evasive_path_match" "config_evasive_path_match" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  enable_path_match = true
}

// USE CASE: user wants to override the evasive path match setting for a security policy
resource "akamai_appsec_advanced_settings_evasive_path_match" "policy_override" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  enable_path_match = true
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Optional) The ID of a specific security policy to which the evasive path match setting should be applied. If not supplied, the indicated setting will be applied to all policies within the configuration.

* `enable_path_match` - (Required) Whether to enable path match.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

