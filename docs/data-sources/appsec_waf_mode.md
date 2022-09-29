---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_waf_mode

**Scopes**: Security policy

Returns information about how the Kona Rule Set rules associated with a security configuration and security policy are updated. The WAF (Web Application Firewall) mode determines whether Kona Rule Sets are automatically updated as part of automated attack groups (`mode = AAG`) or whether you must periodically check for new rules and then manually update those rules yourself (`mode = KRS`).

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/mode](https://techdocs.akamai.com/application-security/reference/get-policy-mode)

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

// USE CASE: User wants to view WAF mode details.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_waf_mode" "waf_mode" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "waf_mode_mode" {
  value = data.akamai_appsec_waf_mode.waf_mode.mode
}
output "waf_mode_current_ruleset" {
  value = data.akamai_appsec_waf_mode.waf_mode.current_ruleset
}
output "waf_mode_eval_status" {
  value = data.akamai_appsec_waf_mode.waf_mode.eval_status
}
output "waf_mode_eval_ruleset" {
  value = data.akamai_appsec_waf_mode.waf_mode.eval_ruleset
}
output "waf_mode_eval_expiration_date" {
  value = data.akamai_appsec_waf_mode.waf_mode.eval_expiration_date
}
output "waf_mode_text" {
  value = data.akamai_appsec_waf_mode.waf_mode.output_text
}
output "waf_mode_json" {
  value = data.akamai_appsec_waf_mode.waf_mode.json
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the Kona Rule Set rules.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the Kona Rule Set rules.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `mode`. Security policy mode, either **KRS** (update manually) or **AAG** (update automatically), For organizations running the Adaptive Security Engine (ASE) beta, you'll get back **ASE_AUTO** for automatic updates or **ASE_MANUAL** for manual updates. Please contact your Akamai representative to learn more about ASE.
- `current_ruleset`. Current ruleset version and the ISO 8601 date the version was introduced.
- `eval_status`. Specifies whether evaluation mode is enabled or disabled.
- `eval_ruleset`. Evaluation ruleset version and the ISO 8601 date the evaluation began.
- `eval_expiration_date`. ISO 8601 timestamp indicating when evaluation mode expires. Valid only if `eval_status` is set to **enabled**.
- `output_text`. Tabular report of the mode information.
- `json`. JSON-formatted list of the mode information.