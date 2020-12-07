---
layout: "akamai"
page_title: "Akamai: Slow Post Protection"
subcategory: "Application Security"
description: |-
 Slow Post Protection
---

# akamai_appsec_slow_post_protection_settings

Use the `akamai_appsec_slow_post_protection_settings` resource to update slow POST protection settings for a specific security configuration version.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user would like to set the slow post protection settings for a given security configuration and version
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_slow_post" "slow_post" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  slow_rate_action = "alert"
  slow_rate_threshold_rate = 10
  slow_rate_threshold_period = 30
  duration_threshold_timeout = 20
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `slow_rate_action` - (Required) Specifies the action that the rule should trigger. Either `alert` or `abort`.

* `slow_rate_threshold_rate` - (Required) The average rate in bytes per second over a period of time that you specify before an action (alert or abort) in the policy triggers.

* `slow_rate_threshold_period` - (Required) The amount of time in seconds of how long the server should accept a request to determine whether a POST request is too slow.

* `duration_threshold_timeout` - (Required) The length of time in seconds by which the edge server must have received the first eight kilobytes of the POST body to avoid triggering the specified action.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the current protection settings.

