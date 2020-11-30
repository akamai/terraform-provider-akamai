---
layout: "akamai"
page_title: "Akamai: WAF Mode"
subcategory: "Application Security"
description: |-
 WAF Mode
---

# akamai_appsec_waf_mode

Use the `akamai_appsec_waf_mode` data source to retrieve the mode that indicates how the WAF rules of the given security configuration version and security policy will be updated.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view waf mode details.
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_waf_mode" "waf_mode" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.policy_id
}

output "waf_mode_mode" {
  value = akamai_appsec_waf_mode.waf_mode.mode
}
output "waf_mode_current_ruleset" {
  value = akamai_appsec_waf_mode.waf_mode.current_ruleset
}
output "waf_mode_eval_status" {
  value = akamai_appsec_waf_mode.waf_mode.eval_status //-- enabled/disabled
}
output "waf_mode_eval_ruleset" {
  value = akamai_appsec_waf_mode.waf_mode.eval_ruleset
}
output "waf_mode_eval_expiration_date" {
  value = akamai_appsec_waf_mode.waf_mode.eval_expiration_date
}
output "waf_mode_text" {
  value = data.akamai_appsec_waf_mode.waf_mode.output_text
}
output "waf_mode_json" {
  value = data.akamai_appsec_waf_mode.waf_mode.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `mode` - The security policy mode, either `KRS` (update manually) or `AAG` (update automatically),

* `current_ruleset` - The current rule set version and the ISO 8601 date the rule set version was introduced; this date acts like a version number. 

* `eval_status` - Whether the evaluation mode is enabled or disabled."

* `eval_ruleset` - The evaluation rule set version and the ISO 8601 date the evaluation starts.

* `eval_expiration_date` - The ISO 8601 time stamp when the evaluation is expiring. This value only appears when `eval` is set to "enabled".

* `output_text` - A tabular display of the mode information.

* `json` - A JSON-formmated list of the mode information.
