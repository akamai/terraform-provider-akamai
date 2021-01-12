---
layout: "akamai"
page_title: "Akamai: SiemSettings"
subcategory: "Application Security"
description: |-
 SiemSettijgs
---

# akamai_appsec_siem_settings

The `akamai_appsec_siem_settings` data source allows you to retrieve the SIEM settings for a specific configuration. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsiemsettings).

## Example Usage

Basic usage:

```hcl
// OPEN API --> https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsiemsettings
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the siem settings with a given security configuration
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_siem_settings" "siem_settings" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
}

output "siem_settings_json" {
  value = data.akamai_appsec_siem_settings.siem_settings.json
}

output "siem_settings_output" {
  value = data.akamai_appsec_siem_settings.siem_settings.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list of the SIEM setting information.

* `output_text` - A tabular display showing the SIEM setting information.

