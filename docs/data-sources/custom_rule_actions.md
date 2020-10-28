---
layout: "akamai"
page_title: "Akamai: CustomRuleActions"
subcategory: "APPSEC"
description: |-
 CustomRuleActions
---

# akamai_appsec_custom_rule_actions

Use the `akamai_appsec_custom_rule_actions` data source to retrieve information about the actions defined for the custom rules associated with a specific security configuration, version and security policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

data "akamai_appsec_custom_rule_actions" "custom_rule_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  policy_id = "crAP_75829"
}

output "custom_rule_actions" {
  value = data.akamai_appsec_custom_rule_actions.custom_rule_actions.output_text
}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `policy_id` - (Required) The ID of the security policy to use

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the ID, name, and action of all custom rules associated with the specified security configuration, version and security policy.

