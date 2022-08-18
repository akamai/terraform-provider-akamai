---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsPragmaHeader"
subcategory: "Application Security"
description: |-
  AdvancedSettingsPragmaHeader
---

# akamai_appsec_advanced_settings_pragma_header

**Scopes**: Security configuration; security policy

Specifies the headers you can exclude from inspection when you are working with a Pragma debug header, a header that provides information about such things as: the edge routers used in a transaction; the Akamai IP addresses involved; whether a request was cached or not; etc. By default, pragma headers are removed from all responses.

This operation can be applied at the security configuration level (in which case it applies to all the security policies in the configuration), or can be customized for an individual security policy.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/pragma-header](https://techdocs.akamai.com/application-security/reference/put-policies-pragma-header)

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

// USE CASE: User wants to configure the pragma header settings for a security policy.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_advanced_settings_pragma_header" "pragma_header" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
  pragma_header      = file("${path.module}/pragma_header.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the pragma header settings being modified.

- `security_policy_id` (Optional). Unique identifier of the security policy associated with the pragma header settings being modified. If not included, pragma header settings are modified at the configuration scope and, as a result, apply to all the security policies associated with the configuration.

- `pragma_header` (Required). Path to a JSON file containing information about the conditions to exclude from the default remove action. By default, the Pragma header debugging information is stripped from an operation's response except in cases where you set `excludeCondition`. 

  To remove existing settings, submit your request with an empty payload ( **{}** ) at the top-level of an object. For example, use the following JSON snippet in the request body to remove the **REQUEST_HEADER_VALUE_MATCH** from the excluded conditions:

  `"type": "{}"`

  Note that, if you submit an empty payload for each member, you'll clear all of your condition settings.

  If you want to modify pragma header settings at the security configuration scope (as opposed to the security policy scope), it's recommended that you first contact your Akamai representative.