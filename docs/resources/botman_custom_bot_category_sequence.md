---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_custom_bot_category_sequence

**Scopes**: Security configuration

Modifies the order in which custom bot categories are evaluated. To set the evaluation order, create a JSON array containing the IDs of your custom bot categories. Categories are evaluated in the same order in which they appear in the array: the category listed first in the array is evaluated first, the category listed second in the array is evaluated second, etc.

Use the [akamai_botman_custom_bot_category_sequence](../data-sources/akamai_botman_custom_bot_category_sequence) data source to review your existing custom category sequence.

**Related API Endpoints**:

- [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-bot-category-sequence](https://techdocs.akamai.com/bot-manager/reference/put-custom-bot-category-sequence). Modifies the custom bot category sequence.

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

resource "akamai_botman_custom_bot_category_sequence" "custom_category_sequence" {
  config_id    = data.akamai_appsec_configuration.configuration.config_id
  category_ids = ["cc9c3f89-e179-4892-89cf-d5e623ba9dc7", "d79285df-e399-43e8-bb0f-c0d980a88e4f", "afa309b8-4fd5-430e-a061-1c61df1d2ac2"]
}
```

To review your current custom bot category sequence, use the [akamai_botman_custom_bot_category_sequence](../data-sources/akamai_botman_custom_bot_category_sequence) data source.

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom bot category sequence.
- `category_ids` (Required). JSON array of custom bot category IDs, with individual IDs separated by using commas. The order of the categories in the array determines the order in which those categories are evaluated.

**See also**:

- [How Bot Manager evaluation works](https://techdocs.akamai.com/bot-manager/docs/how-bot-manager-evaluation-works)
