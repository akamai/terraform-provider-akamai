---
layout: "akamai"
page_title: "Akamai: Rate Policies"
subcategory: "Application Security"
description: |-
 Rate Policies
---

# akamai_appsec_rate_policies

Use the `akamai_appsec_rate_policies` data source to retrieve the rate policies for a specific security configuration, or a single rate policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to see all rate policies associated with a given configuration version
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_rate_policies" "rate_policies" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}
output "rate_policies_output" {
  value = data.akamai_appsec_rate_policies.rate_policies.output_text
}
output "rate_policies_json" {
  value = data.akamai_appsec_rate_policies.rate_policies.json
}

// USE CASE: user wants to see a single rate policy associated with a given Configuration Version
data "akamai_appsec_rate_policies" "rate_policy" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  rate_policy_id = var.rate_policy_id
}
output "rate_policy_json" {
  value = data.akamai_appsec_rate_policies.rate_policy.json
}
output "rate_policy_output" {
  value = data.akamai_appsec_rate_policies.rate_policy.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `rate_policy_id` - (Optional) The ID of the rate policy to use. If this parameter is not supplied, information about all rate policies will be returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the ID and name of all rate policies associated with the specified security configuration.

* `json` - A JSON-formatted list of the rate policy information.

