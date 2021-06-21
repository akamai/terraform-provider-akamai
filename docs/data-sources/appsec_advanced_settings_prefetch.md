---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsPrefetch"
subcategory: "Application Security"
description: |-
 AdvancedSettingsPrefetch
---

# akamai_appsec_advanced_settings_prefetch

Use the `akamai_appsec_advanced_settings_prefetch` data source to retrieve information the prefetch request settings for a security configuration. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getprefetchrequestsforaconfiguration).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the prefetch request settings for a given security configuration
// /appsec/v1/configs/{configId}/versions/{versionNum}/advanced-settings/prefetch
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_advanced_settings_prefetch" "prefetch" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

//tabular data of all fields - 3 boolean fields and one extensions text
output "advanced_settings_prefetch_output" {
  value = data.akamai_appsec_advanced_settings_prefetch.prefetch.output_text
}

output "advanced_settings_prefetch_json" {
  value = data.akamai_appsec_advanced_settings_prefetch.prefetch.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The configuration ID.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list of information about the prefetch request settings.

* `output_text` - A tabular display showing the prefetch request settings.

