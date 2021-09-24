---
layout: "akamai"
page_title: "Akamai: WAF Mode"
subcategory: "Application Security"
description: |-
 WAF Mode
---

# akamai_appsec_waf_mode

**Scopes**: Security policy

Modifies the way your Kona Rule Set rules are updated.
Use **KRS** mode to update the rule sets manually or **AAG** to have those rule sets automatically updated.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/mode](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putmode)

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

// USE CASE: User wants to set the WAF mode.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_waf_mode" "waf_mode" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  mode               = "KRS"
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

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the WAF mode settings being modified.

- `security_policy_id` (Required). Unique identifier of the security policy associated with the WAF mode settings being modified.

- `mode` (Required). Specifies how Kona Rule Set rules are upgraded. Allowed values are:

  - **KRS**. Organizations must manually update their KRS rules.
  - **AAG**. KRS rules are automatically updated by Akamai.
  - **ASE_AUTO**. KRS rules  are automatically updated by Akamai. See the note below for more information.
  - **ASE_MANUAL**. Organizations must manually update their KRS rules. See the note below for more information.

  Note. The **ASE_AUTO** and **ASE_MANUAL** options are only available to organizations running the Adaptive Security Engine (ASE) beta. For more information on ASE, please contact your Akamai representative.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `current_ruleset` – Versioning information for the current Kona Rule Set.
- `eval_ruleset`. Versioning information for the Kona Rule Set being evaluated (if applicable) .
- `eval_status`. Returns **enabled** if an evaluation is currently in progress; otherwise returns **disabled**.
- `eval_expiration_date`. Date on which the evaluation period ends (if applicable).
- `output_text`. Tabular report showing the current rule set, WAF mode and evaluation status.

