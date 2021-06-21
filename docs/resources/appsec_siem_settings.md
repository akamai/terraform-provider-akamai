---
layout: "akamai"
page_title: "Akamai: SIEMSettings"
subcategory: "Application Security"
description: |-
 SIEMSettings
---

# akamai_appsec_siem_settings

Use the `akamai_appsec_siem_settings` resource to mpdate the SIEM integration settings for a specific configuration.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to update the siem settings
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_siem_definitions" "siem_definition" {
  siem_definition_name = var.siem_definition_name
}

data "akamai_appsec_security_policy" "security_policies" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

resource "akamai_appsec_siem_settings" "siem" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  enable_siem = true
  enable_for_all_policies = false
  enable_botman_siem = true
  siem_id = data.akamai_appsec_siem_definitions.siem_definition.id
  security_policy_ids = data.akamai_appsec_security_policy.security_policies.policy_list
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The configuration ID to use.

* `enable_siem` - (Required) Whether you enabled SIEM in a security configuration version.

* `enable_for_all_policies` - (Required) Whether you enabled SIEM for all the security policies in the configuration.

* `enable_botman_siem` - (Required) Whether you enabled SIEM for the Bot Manager events.

* `siem_id` - (Required) An integer that uniquely identifies the SIEM settings.

* `security_policy_ids` - (Required) The list of security policy identifiers for which to enable the SIEM integration.

## Attributes Reference

In addition to the arguments above, the following attribute is exported:

* `output_text` - A tabular display showing the updated SIEM integration settings.

