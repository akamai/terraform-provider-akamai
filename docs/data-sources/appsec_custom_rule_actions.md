---
layout: "akamai"
page_title: "Akamai: Custom Rule Actions"
subcategory: "Application Security"
description: |-
 Custom Rule Actions
---

# akamai_appsec_custom_rule_actions

Use the `akamai_appsec_custom_rule_actions` data source to retrieve information about the actions defined for the custom rules, or a specific custom rule, associated with a specific security configuration version and security policy.

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

* `security_policy_id` - (Required) The ID of the security policy to use

* `custom_rule_id` - (Optional) A specific custom rule for which to retrieve information. If not supplied, information about all custom rules will be returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the ID, name, and action of all custom rules, or of the specific custom rule, associated with the specified security configuration, version and security policy.

