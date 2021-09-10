---
layout: "akamai"
page_title: "Akamai: Rule Upgrade Details"
subcategory: "Application Security"
description: |-
 Rule Upgrade Details
---

# akamai_appsec_rule_upgrade_details

**Scopes**: Security policy

Returns information indicating which of your Kona Rule Sets (if any) need to be updated. A value of **false** indicates that no updates are required.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/rules/upgrade-details](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getupgradedetails)

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

// USE CASE: User wants to view Kona Rule Set upgrade details.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_rule_upgrade_details" "upgrade_details" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}
output "upgrade_details_text" {
  value = data.akamai_appsec_rule_upgrade_details.upgrade_details.output_text
}
output "upgrade_details_json" {
  value = data.akamai_appsec_rule_upgrade_details.upgrade_details.json
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the Kona Rule Sets.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the Kona Rule Sets.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing changes (additions and deletions) to the rules for the specified security policy.
- `json`. JSON-formatted list of the changes (additions and deletions) to the rules for the specified security policy.

