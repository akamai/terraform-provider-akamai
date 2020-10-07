---
layout: "akamai"
page_title: "Akamai: CustomRule"
subcategory: "APPSEC"
description: |-
  CustomRule
---

# resource_akamai_appsec_custom_rule


The `resource_akamai_appsec_custom_rule` resource allows you to create or re-use CustomRules.

If the CustomRule already exists it will be used instead of creating a new one.

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


resource "akamai_appsec_custom_rule" "appseccustomrule" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    custom_rules_json =  file("${path.module}/custom_rules.json")
}
```

## Argument Reference

The following arguments are supported:
"* `config_id`- (Required) The Configuration ID
* `rules` - (Required) Custom Rules File

# Attributes Reference

The following are the return attributes:

*`rule_id` - The RuleID"

