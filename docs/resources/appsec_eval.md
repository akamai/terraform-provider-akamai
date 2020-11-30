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

// OPEN API --> https://developer.akamai.com/api/cloud_security/application_security/v1.html#postevaluationmode

// USE CASE: user wants to set the eval operation
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_eval" "eval_operation" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
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
  value = akamai_appsec_eval.eval_operation.eval_status  // enabled/disabled
}

//TF destroy - stop the eval (i.e eval_action will be STOP)


```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `eval_operation` - (Required) The operation to perform: Start, Stop, Restart, Update, or Complete.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

 * `evaluating_ruleset` - TBD
 * `expiration_date` - TBD
 * `current_ruleset` - TBD
 * `eval_status  // enabled/disabled` - TBD

