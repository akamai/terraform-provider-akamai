---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_bot_category_exception (Beta)

**Scopes**: Security configuration

Updates the bot categories on the category exceptions list.

The exceptions list is  formatted by using JSON and configured as the value for the `bot_category_exception` argument. For information about constructing this JSON file see the links listed in the **Related API Endpoints** section.

To review your current bot category exception list use the [akamai_botman_bot_category_exception](../data-sources/akamai_botman_bot_category_exception) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/transactional-endpoints/bot-protection-exceptions](https://techdocs.akamai.com/bot-manager/reference/put-bot-category-exceptions). Updates the category exceptions list for the specified security policy.

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

resource "akamai_botman_bot_category_exception" "bot_category_exception" {
  config_id              = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id     = "gms1_134637"
  bot_category_exception = file("${path.module}/bot_category_exception.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bot category exceptions list.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the bot category exceptions list.
- `bot_category_exception` (Required). JSON-formatted array of the bot category IDs to be excluded from behavior anomaly detection. In the preceding sample code, the syntax `file("${path.module}/bot_category_exception.json")` points to the location of a JSON file containing the exception list settings and values.

**See also**:

- [Exclude a bot category from behavior anomaly detection](https://techdocs.akamai.com/bot-manager/docs/exclude-bot-category-trans-endpoint)
