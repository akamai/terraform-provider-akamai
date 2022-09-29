---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_slow_post

**Scopes**: Security policy

Returns the slow POST protection settings for the specified security configuration and policy. Slow POST protections help defend a site against attacks that try to tie up the site by using extremely slow requests and responses.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/slow-post](https://techdocs.akamai.com/application-security/reference/get-policy-slow-post)

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

// USE CASE: user wants to view the slow post protection settings associated with a security policy.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_slow_post" "slow_post" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}
output "slow_post_output_text" {
  value = data.akamai_appsec_slow_post.slow_post.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the slow POST settings.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the slow POST settings.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report including the following:
  - **ACTION**. Action taken any time slow POST protection is triggered. Valid values are:
    - **alert**. Record the event.
    - **abort**. Block the request.
  - **SLOW_RATE_THRESHOLD RATE**. Average rate (in bytes per second over the specified time period) allowed before the specified action is triggered.
  - **SLOW_RATE_THRESHOLD PERIOD**. Amount of time (in seconds) that the server should allow a request before marking the request as being too slow.
  - **DURATION_THRESHOLD TIMEOUT**. Maximum amount of time (in seconds) that the first eight kilobytes of the POST body must be received in order to avoid triggering the specified action.