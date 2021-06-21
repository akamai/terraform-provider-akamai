---
layout: "akamai"
page_title: "Akamai: WAF Mode"
subcategory: "Application Security"
description: |-
 WAF Mode
---

# akamai_appsec_waf_mode

Use the `akamai_appsec_waf_mode` resource to specify how your rule sets are updated. Use KRS mode to update the rule sets manually, or AAG to have them update automatically.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to set the waf mode
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_waf_mode" "waf_mode" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
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
  value = akamai_appsec_waf_mode.waf_mode.eval_status
}
output "waf_mode_eval_ruleset" {
  value = akamai_appsec_waf_mode.waf_mode.eval_ruleset
}
output "waf_mode_eval_expiration_date" {
  value = akamai_appsec_waf_mode.waf_mode.eval_expiration_date
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `mode` - (Required) "KRS" to update the rule sets manually, or "AAG" to have them update automatically.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `current_ruleset` - The current rule set.

* `eval_ruleset` - The rule set being evaluated if any.

* `eval_status` - Either `enabled` if an evaluation is currently in progress, or `disabled` otherwise.

* `eval_expiration_date` - The date on which the evaluation period ends.

* `output_text` - A tabular display showing the current rule set, WAF mode and evaluation status (`enabled` if a rule set is currently being evaluated, `disabled` otherwise).

