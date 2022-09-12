---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_akamai_bot_category_action

**Scopes**: Akamai-defined bot category

Modifies the action taken when an Akamai-defined bot category is triggered. 

To review your current bot category actions, use the [data-sources/akamai_botman_akamai_bot_category_action](../data-sources/akamai_botman_akamai_bot_category_action) data source.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/akamai-bot-category-actions/{categoryId}](https://techdocs.akamai.com/bot-manager/reference/put-akamai-bot-category-action). Modifies the action assigned to an Akamai-defined bot category.

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

resource "akamai_botman_akamai_bot_category_action" "bot_category_action" {
  config_id                  = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id         = "gms1_134637"
  category_id                = "2c8add8e-a23c-4c3e-a5c9-8a3dc0d4c0b8"
  akamai_bot_category_action = "file("${path.module}/action.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the bot category.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the bot category.
- `category_id` (Required). Unique identifier of the Akamai bot category being updated.
- `akamai_bot_category_action` (Required). JSON file containing the action to be taken when the bot category is triggered.

**See also**:

- [Predefined actions for bot detections](https://techdocs.akamai.com/bot-manager/docs/predefined-actions-bot)
