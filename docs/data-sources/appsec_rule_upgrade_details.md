---
layout: "akamai"
page_title: "Akamai: Rule Upgrade Details"
subcategory: "Application Security"
description: |-
 Rule Upgrade Details
---

# akamai_appsec_rule_upgrade_details

Use the `akamai_appsec_rule_upgrade_details` data source to retrieve information on changes to the KRS rule sets.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view upgrade details
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_rule_upgrade_details" "upgrade_details" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
}
output "upgrade_details_text" {
  value = data.akamai_appsec_rule_upgrade_details.upgrade_details.output_text
}
output "upgrade_details_json" {
  value = data.akamai_appsec_rule_upgrade_details.upgrade_details.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing changes (additions and deletions) to the rules for the specified security policy.

* `json` - A JSON-formatted list of the changes (additions and deletions) to the rules for the specified security policy.

