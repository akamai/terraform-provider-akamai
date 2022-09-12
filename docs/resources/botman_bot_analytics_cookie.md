---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_bot_analytics_cookie

**Scopes**: Security configuration

Updates setting values for a security configuration’s bot analytics cookie. To configure a cookie settings you need to create a JSON array containing the desired settings and values. That array is then used as the value of the `bot_analytics_cookie` argument.  For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your current bot analytics cookie settings use the [akamai_botman_bot_analytics_cookie](../data-sources/akamai_botman_bot_analytics_cookie) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/bot-analytics-cookie](https://techdocs.akamai.com/bot-manager/reference/put-bot-analytics-cookie-1). Updates the bot analytics cookie settings for the specified security configuration.

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

resource "akamai_botman_bot_analytics_cookie" "bot_analytics_cookie" {
  config_id            = data.akamai_appsec_configuration.configuration.config_id
  bot_analytics_cookie = file("${path.module}/bot_analytics_cookie.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the analytic cookie settings being updated.
- `bot_analytics_cookie` (Required). JSON-formatted collection of analytic cookie settings and their values.  In the preceding sample code, the syntax `file("${path.module}/bot_analytics_cookie.json")` points to the location of a JSON file containing the analytics cookie settings and values.
