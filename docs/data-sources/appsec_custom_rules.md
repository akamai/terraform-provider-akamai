---
layout: "akamai"
page_title: "Akamai: CustomRules"
subcategory: "Application Security"
description: |-
 CustomRules
---


# akamai_appsec_custom_rules

**Scopes**: Security configuration; custom rule

Returns a list of the custom rules defined for a security configuration; you can also use this resource to return information for an individual custom rule. Custom rules are rules you have created yourself and are not part of the Kona Rule Set.

**Related API Endpoint**:[/appsec/v1/configs/{configId}/custom-rules](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getcustomrules)

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

// USE CASE: User wants to view the custom rules associated with a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
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
// USE CASE: User wants to view a specific custom rule.

data "akamai_appsec_custom_rules" "specific_custom_rule" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  custom_rule_id = "60029316"
}
output "specific_custom_rule_json" {
  value = data.akamai_appsec_custom_rules.specific_custom_rule.json
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom rules
- `custom_rule_id` (Optional). Unique identifier of the custom rule you want to return information for. If not included, information is returned for all your custom rules.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the ID and name of the custom rule information.
- `json`. JSON-formatted report of the custom rule information.

