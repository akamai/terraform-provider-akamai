---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsPrefetch"
subcategory: "Application Security"
description: |-
 AdvancedSettingsPrefetch
---

# akamai_appsec_advanced_settings_prefetch

**Scopes**: Security configuration

Returns information about your prefetch request settings. By default, Web Application Firewall inspects only external requests — requests originating outside of your firewall or Akamai's edge servers. When prefetch is enabled, requests between your origin servers and Akamai's edge servers can also be inspected by the firewall. The returned information is described in the [PrefetchRequest members](https://developer.akamai.com/api/cloud_security/application_security/v1.html#deb7220d) section of the Application Security API.

**Related** **API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/prefetch](https://techdocs.akamai.com/application-security/reference/get-advanced-settings-prefetch)

## Example usage

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

// USE CASE: User wants to view the prefetch request settings for a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_advanced_settings_prefetch" "prefetch" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

//USE CASE: User wants to display returned data in a table.

output "advanced_settings_prefetch_output" {
  value = data.akamai_appsec_advanced_settings_prefetch.prefetch.output_text
}

output "advanced_settings_prefetch_json" {
  value = data.akamai_appsec_advanced_settings_prefetch.prefetch.json
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the prefetch settings.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of information about the prefetch request settings.
- `output_text`. Tabular report showing the prefetch request settings.