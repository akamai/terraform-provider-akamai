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
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

data "akamai_appsec_custom_rules" "custom_rules" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "custom_rules_list" {
  value = data.akamai_appsec_custom_rules.custom_rules.output_text
}
```

## Argument Reference

The following argument is supported:

* `config_id` - (Required) The ID of the security configuration to use.

## Attributes Reference

In addition to the argument above, the following attribute is exported:

* `output_text` - A tabular display showing the ID and name of the custom rules defined for the security configuration.

