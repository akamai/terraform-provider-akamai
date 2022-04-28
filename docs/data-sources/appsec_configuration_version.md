---
layout: "akamai"
page_title: "Akamai: ConfigurationVersion"
subcategory: "Application Security"
description: |-
 ConfigurationVersion
---


# akamai_appsec_configuration_version

**Scopes**: Security configuration

Returns versioning information for a security configuration.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}](https://techdocs.akamai.com/application-security/reference/get-version-number)

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

data "akamai_appsec_configuration" "specific_configuration" {
  name = "Documentation"
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
  version   = 42
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

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration you want to return version information for.
- `version` (Optional). Version number of the security configuration you want to return information about. If not included, information about all the security configuration's versions is returned.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `latest_version`. Most-recent version of the security configuration.

- `staging_status`. Status of the specified version in staging. Valid values are:

  - **Active**
  - **Inactive**
  - **Deactivated**

  Returned only if the `version` argument is included in the Terraform configuration file.

- `production_status`. Status of the specified version in production. Valid values are:

  - **Active**
  - **Inactive**
  - **Deactivated**

  Returned only if the `version` argument is included in the Terraform configuration file.

- `output_text`. Tabular report showing the version number, staging status, and production status properties and property values.