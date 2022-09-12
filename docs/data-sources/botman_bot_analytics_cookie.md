---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_bot_analytics_cookie

**Scopes**: Security configuration

Returns information about the bot analytics cookie used by the specified security configuration.

Use the [akamai_botman_bot_analytics_cookie](../resources/akamai_botman_bot_analytics_cookie) resource to modify your bot analytics cookie.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/bot-analytics-cookie](https://techdocs.akamai.com/bot-manager/reference/get-bot-analytics-cookie-1). Returns settings information for the bot analytics cookie.

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

data "akamai_botman_bot_analytics_cookie" "analytics_cookies" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "analytics_cookies_json" {
  value = data.akamai_botman_bot_analytics_cookie.analytic_cookies.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bot analytics cookie.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your bot analytic cookie settings, including the cookie name and the hostnames associated with the cookie.
