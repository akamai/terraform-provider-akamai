---
layout: "akamai"
page_title: "Akamai: CustomDeny"
subcategory: "Application Security"
description: |-
  CustomDeny
---

# akamai_appsec_custom_deny

**Scopes**: Custom deny

Modifies a custom deny action. Custom denies enable you to craft your own error message or redirect pages for use when HTTP requests are denied.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/custom-deny](https://techdocs.akamai.com/application-security/reference/get-custom-deny-actions)

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

// USE CASE: User wants to create a custom deny action by using a JSON-formatted definition.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_custom_deny" "custom_deny" {
  config_id   = data.akamai_appsec_configuration.configuration.config_id
  custom_deny = file("${path.module}/custom_deny.json")
}

output "custom_deny_id" {
  value = akamai_appsec_custom_deny.custom_deny.custom_deny_id
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the custom deny.
- `custom_deny` (Required). Path to a JSON file containing properties and property values for the custom deny. 

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `custom_deny_id`. ID of the new custom deny action.