---
layout: "akamai"
page_title: "Akamai: Rate Policy"
subcategory: "Application Security"
description: |-
  Rate Policy
---

# akamai_appsec_rate_policy

**Scopes**: Security configuration; rate policy

Creates, modifies, or deletes rate policies. Rate polices help you monitor and moderate the number and rate of all the requests you receive.
In turn, this helps you prevent your website from being overwhelmed by a dramatic and unexpected surge in traffic.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/rate-policies](https://techdocs.akamai.com/application-security/reference/post-rate-policies)

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

// USE CASE: User wants to create a rate policy for a security configuration by using a JSON-formatted rule definition.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_rate_policy" "rate_policy" {
  config_id   = data.akamai_appsec_configuration.configuration.config_id
  rate_policy = file("${path.module}/rate_policy.json")
}
output "rate_policy_id" {
  value = akamai_appsec_rate_policy.rate_policy.rate_policy_id
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the rate policy being modified.
- `rate_policy` (Required). Path to a JSON file containing a rate policy definition. You can view a sample rate policy JSON file in the [RatePolicy](https://developer.akamai.com/api/cloud_security/application_security/v1.html#ratepolicy) section of the Application Security API documentation.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `rate_policy_id`. ID of the modified or newly-created rate policy.
