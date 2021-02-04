---
layout: "akamai"
page_title: "Akamai: CustomRules"
subcategory: "Application Security"
description: |-
 CustomRules
---

# akamai_appsec_custom_rules

Use the `akamai_appsec_custom_rules` data source to retrieve a list of the custom rules defined for a security configuration.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to see the custom rules associated with a given security configuration
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_custom_rules" "custom_rules" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}
output "custom_rules_output_text" {
  value = data.akamai_appsec_custom_rules.custom_rules.output_text
}
output "custom_rules_json" {
  value = data.akamai_appsec_custom_rules.custom_rules.json
}
output "custom_rules_config_id" {
  value = data.akamai_appsec_custom_rules.custom_rules.config_id
}
// USE CASE: user wants to see a specific custom rule
data "akamai_appsec_custom_rules" "specific_custom_rule" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  custom_rule_id = var.custom_rule_id
}
output "specific_custom_rule_json" {
  value = data.akamai_appsec_custom_rules.specific_custom_rule.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `custom_rule_id` - (Optional) The ID of a specific custom rule to use. If not supplied, information about all custom rules associated with the given security configuration will be returned.

## Attributes Reference

In addition to the argument above, the following attribute is exported:

* `output_text` - A tabular display showing the ID and name of the custom rule(s).

* `json` - A JSON-formatted display of the custom rule information.

