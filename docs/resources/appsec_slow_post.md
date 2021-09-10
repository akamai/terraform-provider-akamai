---
layout: "akamai"
page_title: "Akamai: Slow Post"
subcategory: "Application Security"
description: |-
 Slow Post
---

# akamai_appsec_slow_post

**Scopes**: Security policy

Modifies slow POST protection settings for a security configuration and security policy. Slow POST protections help defend a site against attacks that try to tie up the site by using extremely slow requests and responses.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/slow-post](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putslowpostprotectionsettings)

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

// USE CASE: User wants to set slow post protection settings for a security configuration and security policy.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_slow_post" "slow_post" {
  config_id                  = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id         = "gms1_134637"
  slow_rate_action           = "alert"
  slow_rate_threshold_rate   = 10
  slow_rate_threshold_period = 30
  duration_threshold_timeout = 20
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the slow POST settings being modified.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the slow POST settings being modified.
- `slow_rate_action` (Required). Action to be taken if slow POST protection is triggered. Allowed values are:
  - **alert**. Record the event.
  - **abort**. Block the request.
- `slow_rate_threshold_rate` (Optional). Average rate (in bytes per second over the specified time period) allowed before the specified action is triggered.
- `slow_rate_threshold_period` (Optional). Amount of time (in seconds) that the server should allow a request before marking the request as being too slow.
- `duration_threshold_timeout` (Optional). Maximum amount of time (in seconds) that the first eight kilobytes of the POST body must be received in to avoid triggering the specified action.

