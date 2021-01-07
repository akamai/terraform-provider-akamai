---
layout: "akamai"
page_title: "Akamai: KRS Rule Condition-Exception"
subcategory: "Application Security"
description: |-
 KRS Rule Condition-Exception
---

# akamai_appsec_rule_condition_exception

Use the `akamai_appsec_rule_condition_exception` data source to retrieve a KRS rule’s conditions and exceptions.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view condition-exception for a rule
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_rule_condition_exception" "condition_exception" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.policy_id
  rule_id = var.rule_id
}
output "condition_exception_text" {
  value = data.akamai_appsec_rule_condition_exception.condition_exception.output_text
}

output "condition_exception_json" {
  value = data.akamai_appsec_rule_condition_exception.condition_exception.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `rule_id` - (Required) The ID of the rule to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing boolean values indicating whether conditions and exceptions are present

* `json` - A JSON-formatted list of the condition and exception information for the specified rule.

