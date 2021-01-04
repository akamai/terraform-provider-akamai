---
layout: "akamai"
page_title: "Akamai: Custom Rule"
subcategory: "Application Security"
description: |-
  Custom Rule
---

# akamai_appsec_custom_rule


The `akamai_appsec_custom_rule` resource allows you to create or modify a custom rule associated with a given security configuration.

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
  custom_rule = data.local_file.rules.content
}

output "custom_rule_rule_id" {
  value = akamai_appsec_custom_rule.custom_rule.custom_rule_id
}
```

## Argument Reference

* `config_id` - (Required) The ID of the security configuration to use.

* `custom_rule` - (Required) The name of a JSON file containing a custom rule definition ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#postcustomrules)).


## Attribute Reference

In addition to the arguments above, the following attribute is exported:

* `custom_rule_id` - The ID of the custom rule.

