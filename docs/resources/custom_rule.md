---
layout: "akamai"
page_title: "Akamai: CustomRule"
subcategory: "APPSEC"
description: |-
  CustomRule
---

# resource_akamai_appsec_custom_rule


The `resource_akamai_appsec_custom_rule` resource allows you to create or modify a custom rule associated with a given security configuration.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

data "local_file" "rules" {
  filename = "${path.module}/custom_rules.json"
}

resource "akamai_appsec_custom_rule" "custom_rule" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  rules = data.local_file.rules.content
}

output "custom_rule_rule_id" {
  value = akamai_appsec_custom_rule.custom_rule.rule_id
}
```

## Argument Reference

* `config_id` - (Required) The ID of the security configuration to use.

* `rules` - (Required) The name of a JSON file containing a custom rule definition.


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `rule_id` - The ID of the custom rule.

