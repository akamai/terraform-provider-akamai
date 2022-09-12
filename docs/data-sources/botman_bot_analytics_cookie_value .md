---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_bot_analytics_cookie_value

**Scopes**: Universal (all bot analytic cookie values)

Returns information about the bots assigned to your bot analytics cookie value. 

**Related API Endpoints**:

- [/appsec/v1/bot-analytics-cookie/values](https://techdocs.akamai.com/bot-manager/reference/get-bot-analytics-cookie-values). Returns all the bots associated with your analytics cookie.

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

data "akamai_botman_bot_analytics_cookie_value" "analytics_cookie_values" {
}

output "analytics_cookie_values_json" {
  value = data.akamai_botman_bot_analytics_cookie_value.analytics_cookie_values.json
}
```

## Argument Reference

This data source doesn’t accept any arguments.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing the IDs of all the bots associated with the analytics cookie.

**See also**:

- [Forward bot status to your analytics tool](https://techdocs.akamai.com/bot-manager/docs/forward-bot-status-to-your-analytics-tool)
