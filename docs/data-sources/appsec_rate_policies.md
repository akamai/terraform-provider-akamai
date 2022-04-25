---
layout: "akamai"
page_title: "Akamai: Rate Policies"
subcategory: "Application Security"
description: |-
 Rate Policies
---

# akamai_appsec_rate_policies

**Scopes**: Security configuration; rate policy

Returns information about your rate policies. Rate polices help you monitor and moderate the number and rate of all the requests you receive; in turn, this helps you prevent your website from being overwhelmed by a dramatic, and unexpected, surge in traffic.

**Related API Endpoint:** [/appsec/v1/configs/{configId}/versions/{versionNumber}/rate-policies](https://techdocs.akamai.com/application-security/reference/get-rate-policies)

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

// USE CASE: User wants to view all the rate policies associated with a configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
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

// USE CASE: User wants to see a specific rate policy.

data "akamai_appsec_rate_policies" "rate_policy" {
  config_id      = data.akamai_appsec_configuration.configuration.config_id
  rate_policy_id = "122149"
}
output "rate_policy_json" {
  value = data.akamai_appsec_rate_policies.rate_policy.json
}
output "rate_policy_output" {
  value = data.akamai_appsec_rate_policies.rate_policy.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the rate policies.
- `rate_policy_id` (Optional). Unique identifier of the rate policy you want to return information for. If not included, information is returned for all your rate policies.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the ID and name of the rate policies.
- `json`. JSON-formatted list of the rate policy information.