---
layout: "akamai"
page_title: "Akamai: Evaluation"
subcategory: "Application Security"
description: |-
 Evaluation
---

# akamai_appsec_eval

**Scopes**: Security policy

Issues an evaluation mode command (`Start`, `Stop`, `Restart`, `Update`, or `Complete`) to a security configuration.
Evaluation mode is used for testing and fine-tuning your Kona Rule Set rules and configuration settings.
In evaluation mode rules are triggered by events, but the only thing those rules do is record the actions they *would* have taken had the event occurred on the production network.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/eval](https://techdocs.akamai.com/application-security/reference/post-policy-eval)

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

// USE CASE: User wants to issue an evaluation mode command.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_eval" "eval_operation" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  eval_operation     = "START"
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

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration where evaluation mode will take place (or is currently taking place).
- `security_policy_id` (Required). Unique identifier of the security policy associated with the evaluation process.
- `eval_operation` (Required). Evaluation mode operation. Allowed values are:
  - **START**. Starts evaluation mode. By default, evaluation mode runs for four weeks.
  - **STOP**, Pauses evaluation mode without upgrading the Kona Rule Set on your production network.
  - **RESTART**. Resumes an evaluation trial that was paused by using the **STOP** command.
  - **UPDATE**. Upgrades the Kona Rule Set rules in the evaluation ruleset to their latest versions.
  - **COMPLETE**. Concludes the evaluation period (even if the four-week trial mode is not over) and automatically upgrades the Kona Rule Set on your production network to the same rule set you just finished evaluating.
- `eval_mode` (Optional). Set to **ASE_AUTO** to have your Kona Rule Set rules automatically updated during the evaluation period; set to **ASE_MANUAL** if you want to manually update your evaluation rules. Note that this option is only available to organizations running the Adaptive Security Engine (ASE) beta. For more information about ASE, please contact your Akamai representative.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `evaluating_ruleset`. Versioning information for the Kona Rule Set being evaluated.
- `expiration_date`. Date when the evaluation period ends.
- `current_ruleset`. Versioning information for the Kona Rule Set currently in use on the production network.
- `eval_status`. If **true**, an evaluation is currently in progress; if **false**, evaluation is either paused or is not running.
