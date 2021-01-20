---
layout: "akamai"
page_title: "Akamai: Attack Group Condition Exception"
subcategory: "Application Security"
description: |-
 Attack Group Condition Exception
---

# akamai_appsec_attack_group_condition_exception

Use the `akamai_appsec_attack_group_condition_exception` data source to retrieve an attack group's conditions and exceptions.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view condition-exception for an attack group
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_attack_group_condition_exception" "condition_exception" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  attack_group = var.attack_group
}
output "condition_exception_text" {
  value = data.akamai_appsec_attack_group_condition_exception.condition_exception.output_text
}
output "condition_exception_json" {
  value = data.akamai_appsec_attack_group_condition_exception.condition_exception.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `attack_group` - (Required) The attack group to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the condition and exception information.

* `json` - The condition and exception information in JSON format.

