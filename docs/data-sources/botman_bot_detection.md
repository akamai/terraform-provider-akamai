---
layout: akamai
subcategory: Bot Manager
---

# akamai_botman_bot_detection (Beta)

**Scopes**: Universal (all bot detection methods); bot detection method

Returns information about all available bot detection methods. Bot detection methods are used to help determine whether a request came from a bot or from a human being.

Use the `detection_name` argument to return data only for the specified detection method.

**Related API Endpoints**:

- [/appsec/v1/bot-detections](https://techdocs.akamai.com/bot-manager/reference/get-bot-detections). Returns information about all the available bot detection methods.
- [/appsec/v1/bot-detections/{detectionId}](https://techdocs.akamai.com/bot-manager/reference/get-bot-detection). Returns information for the specified bot detection method.

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

// USE CASE: User wants to return information for all bot detection methods

data "akamai_botman_bot_detection" "detection_methods" {
}

// USE CASE: User only wants to return information fpr the Impersonators of Known Bots detection method

data "akamai_botman_bot_detection" "detection_method" {
  detection_name = "Impersonators of Known Bots "
}

output "detection_method_json" {
  value = data.akamai_botman_bot_detection.detection_method.json
}
```

## Argument Reference

This resource supports the following arguments:

- `detection_name` (Optional). Unique name of a bot detection method.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `json`. JSON-formatted output containing information about your bot detection methods. The returned information includes a Boolean field indicating whether the detection method is active or inactive.

**See also**:

- [Bot detection methods and rule IDs](https://techdocs.akamai.com/bot-manager/docs/bot-det-methods-rule-ids)
