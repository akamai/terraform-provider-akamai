---
layout: "akamai"
page_title: "Akamai: Rate Policy Actions"
subcategory: "Application Security"
description: |-
 Rate Policy Actions
---

# akamai_appsec_rate_policy_actions

Use the `akamai_appsec_rate_policy_actions` data source to retrieve a list of all rate policies associated with a given configuration version and security policy, or the actions associated with a specific rate policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view the all rate policy actions associated with a given security policy
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_rate_policy_actions" "rate_policy_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
}
output "rate_policy_actions" {
  value = data.akamai_appsec_rate_policy_actions.rate_policy_actions.output_text
}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `rate_policy_id` - (Optional) The ID of the rate policy to use. If not supplied, information about all rate policies will be returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the ID IPv4Action and IPv6Action of the indicated security policy or policies.

