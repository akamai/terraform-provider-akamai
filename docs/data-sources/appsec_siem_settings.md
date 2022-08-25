---
layout: "akamai"
page_title: "Akamai: SiemSettings"
subcategory: "Application Security"
description: |-
 SiemSettings
---

# akamai_appsec_siem_settings

**Scopes**: Security configuration

Returns the SIEM (Security Event and Information Management) settings for a security configuration. 

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/siem](https://techdocs.akamai.com/application-security/reference/get-siem)

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

// USE CASE: User wants to view the SIEM settings for a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_siem_settings" "siem_settings" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "siem_settings_json" {
  value = data.akamai_appsec_siem_settings.siem_settings.json
}

output "siem_settings_output" {
  value = data.akamai_appsec_siem_settings.siem_settings.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration you want to return information for.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of the SIEM setting information.
- `output_text`. Tabular report showing the SIEM setting information.