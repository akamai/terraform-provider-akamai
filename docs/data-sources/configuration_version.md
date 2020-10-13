---
layout: "akamai"
page_title: "Akamai: ConfigurationVersion"
subcategory: "APPSEC"
description: |-
 ConfigurationVersion
---

# akamai_appsec_configuration_version

Use the `akamai_appsec_configuration_version` data source to retrieve information about the versions of a security configuration, or about a specific version.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "specific_configuration" {
  name = "Akamai Tools"
}

data "akamai_appsec_configuration_version" "versions" {
  config_id = data.akamai_appsec_configuration.specific_configuration.config_id
}

output "versions_output_text" {
  value = data.akamai_appsec_configuration_version.versions.output_text
}

output "versions_latest" {
  value = data.akamai_appsec_configuration_version.versions.latest_version
}

data "akamai_appsec_configuration_version" "specific_version" {
  config_id = data.akamai_appsec_configuration.specific_configuration.config_id
  version = 42
}

output "specific_version_version" {
  value = data.akamai_appsec_configuration_version.specific_version.version
}

output "specific_version_staging" {
  value = data.akamai_appsec_configuration_version.specific_version.staging_status
}

output "specific_version_production" {
  value = data.akamai_appsec_configuration_version.specific_version.production_status
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Optional) The version number of the security configuration to use. If not supplied, information about all versions of the specified security configuration is returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `latest_version` - The last version of the security configuration created.

* `staging_status` - The status of the specified version in staging: "Active", "Inactive", or "Deactivated". Returned only if `version` was specified.

* `production_status` - The status of the specified version in production: "Active", "Inactive", or "Deactivated". Returned only if `version` was specified.

* `output_text` - A tabular listing showing the following information about all versions of the security configuration: version number, staging status, and production status.

