---
layout: "akamai"
page_title: "Akamai: Mode"
subcategory: "Application Security"
description: |-
 Mode
---

# akamai_appsec_mode

Use the `akamai_appsec_mode` resource to specify how your rule sets are updated. Use KRS mode to update the rule sets manually, or AAG to have them update automatically.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// OPEN API --> https://developer.akamai.com/api/cloud_security/application_security/v1.html#putmode

// USE CASE: user wants to set the waf mode
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_waf_mode" "waf_mode" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.policy_id
  mode = var.mode
}

output "waf_mode_mode" {
  value = akamai_appsec_waf_mode.waf_mode.mode
}

output "waf_mode_current_ruleset" {
  value = akamai_appsec_waf_mode.waf_mode.current_ruleset
}

output "waf_mode_eval_status" {
  value = akamai_appsec_waf_mode.waf_mode.eval_status  // enabled/disabled
}

output "waf_mode_eval_ruleset" {
  value = akamai_appsec_waf_mode.waf_mode.eval_ruleset
}

output "waf_mode_eval_expiration_date" {
  value = akamai_appsec_waf_mode.waf_mode.eval_expiration_date
}

//TF destroy - no-op


```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `mode` - (Required) "KRS" to update the rule sets manually, or "AAG" to have them update automatically.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - TBD


