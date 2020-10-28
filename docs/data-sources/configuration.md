---
layout: "akamai"
page_title: "Akamai: Configuration"
subcategory: "APPSEC"
description: |-
 Configuration
---

# akamai_appsec_configuration

Use the `akamai_appsec_configuration` data source to retrieve the list of security configurations, or information about a specific security configuration.


## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configurations" {
}

output "configuration_list" {
  value = data.akamai_appsec_configuration.configurations.output_text
}

data "akamai_appsec_configuration" "specific_configuration" {
  name = "Akamai Tools"
}

output "latest" {
  value = data.akamai_appsec_configuration.specific_configuration.latest_version
}

output "staging" {
  value = data.akamai_appsec_configuration.specific_configuration.staging_version
}

output "production" {
  value = data.akamai_appsec_configuration.specific_configuration.production_version
}

output "id" {
  value = data.akamai_appsec_configuration.specific_configuration.config_id
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of a specific security configuration. If not supplied, information about all security configurations is returned.

* `version` - (Optional) The specific version number to return. If specified, this value is returned for use in specifying other data sources or resources.


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `config_id` - The ID of the specified security configuration. Returned only if `name` was specified.

* `output_text` - A tabular display showing the following information about all available security configurations: config_id, name, latest version, version active in staging, and version active in production.

* `latest_version` - The last version of the specified security configuration created. Returned only if `name` was specified.

* `staging_version` - The version of the specified security configuration currently active in staging. Returned only if `name` was specified.

* `production_version` - The version of the specified security configuration currently active in production. Returned only if `name` was specified.
