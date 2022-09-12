---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_bot_category_exception

**Scopes**: Security policy

Returns information about all the bots that have been excluded from transactional endpoint bot protection. 

Use the [akamai_botman_bot_category_exception](../resources/akamai_botman_bot_category_exception) resource to update your bot category exception list.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/transactional-endpoints/bot-protection-exceptions](https://techdocs.akamai.com/bot-manager/reference/get-bot-category-exceptions). Returns information about all the bots excluded from behavior anomaly detection.

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

data "akamai_botman_bot_category_exception" "category_exception" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
}

output "category_exception_json" {
  value = data.akamai_botman_bot_category_exception.category_exception.json
}
```

# Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bot category exceptions.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the bot category exceptions.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about the bot IDs (including bots in both Akamai-defined categories and in custom categories) that have been excluded from behavior anomaly detection.

**See also**:

- [Exclude a bot category from behavior anomaly detection](https://techdocs.akamai.com/bot-manager/docs/exclude-bot-category-trans-endpoint)
