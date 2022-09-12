---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_custom_bot_category_sequence

**Scopes**: Security configuration

Returns the category sequence for your custom bot categories. The category sequence determines the order in which custom bot categories are evaluated.

Use the [akamai_botman_custom_bot_category_sequence](../resources/akamai_botman_custom_bot_category_sequence) resource to modify your custom bot category sequence.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-bot-category-sequence](https://techdocs.akamai.com/bot-manager/reference/get-custom-bot-category-sequence). Returns the order in which custom bot categories are evaluated.

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

data "akamai_botman_custom_bot_category_sequence" "category_sequence" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "category_sequence_json" {
  value = data.akamai_botman_custom_bot_category_sequence.category_sequence.json
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom bot category sequence.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing the IDs of your custom bot categories and the order in which those categories are evaluated.

**See also**:

- [How Bot Manager evaluation works](https://techdocs.akamai.com/bot-manager/docs/how-bot-manager-evaluation-works)
