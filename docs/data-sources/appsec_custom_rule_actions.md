---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_custom_rule_actions

**Scopes**: Security policy; custom rule

Retrieve information about the actions defined for your custom rules. Custom rules are rules that you create yourself — these rules aren't part of Akamai's Kona Rule Set.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/custom-rules](https://techdocs.akamai.com/application-security/reference/get-custom-rules)

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
data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_custom_rule_actions" "custom_rule_actions" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}
output "custom_rule_actions" {
  value = data.akamai_appsec_custom_rule_actions.custom_rule_actions.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom rules.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the custom rules.
- `custom_rule_id` (Optional). Unique identifier of the custom rule you want to return information for. If not included, action information is returned for all your custom rules.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the ID, name, and action of the custom rules.