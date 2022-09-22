---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_advanced_settings_prefetch

**Scopes**: Security configuration

Enables inspection of internal requests (that is, requests between your origin servers and Akamai's edge servers). You can also use this resource to apply rate controls to prefetch requests.

When prefetch is enabled, internal requests are inspected by your firewall the same way that external requests (requests that originate outside the firewall and outside Akamai's edge servers) are inspected.

This operation applies at the security configuration level, meaning that the settings affect all the security policies in that configuration.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/advanced-settings/prefetch](https://techdocs.akamai.com/application-security/reference/put-advanced-settings-prefetch)

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

// USE CASE: User wants to configure prefetch settings.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_advanced_settings_prefetch" "prefetch" {
  config_id            = data.akamai_appsec_configuration.configuration.config_id
  enable_app_layer     = false
  all_extensions       = true
  enable_rate_controls = false
  extensions           = [".tiff", ".bmp", ".jpg", ".gif", ".png"]
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the prefetch settings being modified.
- `enable_app_layer` (Required). Set to **true** to enable prefetch requests; set to **false** to disable prefetch requests.
- `all_extensions` (Required). Set to **true** to enable prefetch requests for all file extensions; set to **false** to enable prefetch requests on only a specified set of file extensions. If set to false you must include the `extensions` argument.
- `enable_rate_controls` (Required). Set to **true** to enable prefetch requests for rate controls; set to **false** to disable prefetch requests for rate controls.
- `extensions` (Required). If `all_extensions` is **false**, this must be a JSON array of all the file extensions for which prefetch requests are enabled: prefetch requests won't be used with any file extensions not included in the array. If `all_extensions` is **true**, then this argument must be set to an empty array: **[]**.