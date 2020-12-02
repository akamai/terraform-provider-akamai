---
layout: "akamai"
page_title: "Akamai: Slow Post"
subcategory: "Application Security"
description: |-
 Slow Post
---

# akamai_appsec_slow_post

Use the `akamai_appsec_slow_post` data source to retrieve the slow post protection settings for a given security configuration version and policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to see the slow post protection settings associated with a given
//           security configuration, version and security policy
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_slow_post" "slow_post" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
}
output "slow_post_output_text" {
  value = data.akamai_appsec_slow_post.slow_post.output_text
}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display including the following columns:
  * `ACTION` - The action that the rule should trigger (either `alert` or `abort`)
  * `SLOW_RATE_THRESHOLD RATE` - The average rate in bytes per second over the period specified by `period` before the specified `action` is triggered.
  * `SLOW_RATE_THRESHOLD PERIOD` - The length in seconds of the period during which the server should accept a request before determining whether a POST request is too slow.
  * `DURATION_THRESHOLD TIMEOUT` - The time in seconds before the first eight kilobytes of the POST body must be received to avoid triggering the specified `action`.

