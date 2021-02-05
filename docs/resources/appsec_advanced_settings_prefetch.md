---
layout: "akamai"
page_title: "Akamai: AdvancedSettingsPrefetch"
subcategory: "Application Security"
description: |-
  AdvancedSettingsPrefetch
---

# resource_akamai_appsec_advanced_settings_prefetch

The `resource_akamai_appsec_advanced_settings_prefetch` resource allows you to enable inspection of internal requests (those between your origin and Akamai’s servers) for file types that you specify. You can also apply rate controls to prefetch requests. This operation applies at the configuration level.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to set the prefetch settings
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

resource "akamai_appsec_advanced_settings_prefetch" "prefetch" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  enable_app_layer = false
  all_extensions = true
  enable_rate_controls = false
  extensions = var.extensions
}

output "prefetch_settings" {
  value = akamai_appsec_advanced_settings_prefetch.prefetch.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `enable_app_layer` - (Required) Whether to enable prefetch requests.

* `all_extensions` - (Required) Whether to enable prefetch requests for all extensions.

* `enable_rate_controls` - (Required) Whether to enable prefetch requests for rate controls.

* `extensions` - (Required) The specific extensions for which to enable prefetch requests. If `all_extensions` is True, `extensions` must be an empty list.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the following updated prefetch settings.
  * `ENABLE APP LAYER`
  * `ALL EXTENSION`
  * `ENABLE RATE CONTROLS`
  * `EXTENSIONS`
