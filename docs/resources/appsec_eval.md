---
layout: "akamai"
page_title: "Akamai: Evaluation"
subcategory: "Application Security"
description: |-
 Evaluation
---

# akamai_appsec_eval

Use the `akamai_appsec_eval` resource to perform evaluation mode operations such as Start, Stop, Restart, Update, or Complete.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to set the eval operation
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_eval" "eval_operation" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  eval_operation = var.eval_operation
}
output "eval_mode_evaluating_ruleset" {
  value = akamai_appsec_eval.eval_operation.evaluating_ruleset
}
output "eval_mode_expiration_date" {
  value = akamai_appsec_eval.eval_operation.expiration_date
}
output "eval_mode_current_ruleset" {
  value = akamai_appsec_eval.eval_operation.current_ruleset
}
output "eval_mode_status" {
  value = akamai_appsec_eval.eval_operation.eval_status
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `eval_operation` - (Required) The operation to perform: START, STOP, RESTART, UPDATE, or COMPLETE.

* `eval_mode` - __ASE Beta__. (Optional) Used for ASE Rulesets: ASE_MANUAL or ASE_AUTO - default. "ASE (Adaptive Security Engine) is currently in beta. Please contact your Akamai representative to learn more. Policy Evaluation Rule Actions and Threat Intelligence setting are read only in ASE_AUTO evaluation mode

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `evaluating_ruleset` - The set of rules being evaluated.

* `expiration_date` - The date on which the evaluation period ends.

* `current_ruleset` - The set of rules currently in effect.

* `eval_status` - Either `enabled` if an evaluation is currently in progress (that is, if the `eval_operation` parameter was `START`, `RESTART`, or `COMPLETE`) or `disabled` otherwise (that is, if the `eval_operation` parameter was `STOP` or `UPDATE`).

