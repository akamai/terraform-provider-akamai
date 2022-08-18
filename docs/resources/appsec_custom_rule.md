---
layout: "akamai"
page_title: "Akamai: Custom Rule"
subcategory: "Application Security"
description: |-
  Custom Rule
---

# akamai_appsec_custom_rule

**Scopes**: Security configuration

Creates a custom rule associated with a security configuration. Custom rules are rules that you define yourself and are not part of the Kona Rule Set.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/custom-rules]https://techdocs.akamai.com/application-security/reference/get-configs-custom-rules)

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

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "local_file" "rules" {
  filename = "${path.module}/custom_rules.json"
}

resource "akamai_appsec_custom_rule" "custom_rule" {
  config_id   = data.akamai_appsec_configuration.configuration.config_id
  custom_rule = data.local_file.rules.content
}

output "custom_rule_rule_id" {
  value = akamai_appsec_custom_rule.custom_rule.custom_rule_id
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom rule being modified.
- `custom_rule` (Required). Path to a JSON file containing the custom rule definition. To view a sample JSON file, see the [Create a custom rule](https://techdocs.akamai.com/application-security/reference/post-config-custom-rules) section of the Application Security API documentation.

## Attribute Reference

In addition to the arguments above, the following attribute is exported:

- `custom_rule_id`. ID of the new custom rule.