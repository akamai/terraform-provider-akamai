---
layout: "akamai"
page_title: "Akamai: CustomRules"
subcategory: "APPSEC"
description: |-
 CustomRules
---

# akamai_appsec_custom_rules

Use `akamai_appsec_custom_rules` data source to retrieve a custom_rules id.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}
data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Example for EDGE"
  
}



output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}


data "akamai_appsec_custom_rules" "appseccustomrule" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
}

output "appseccustomrules" {
  value = data.akamai_appsec_custom_rules.appseccustomrule.output_text
}

```

## Argument Reference

The following arguments are supported:

"* `config_id`- (Required) The Configuration ID
* `rules` - (Required) Custom Rules File

# Attributes Reference

The following are the return attributes:

*`output_text` - rule id and name"

