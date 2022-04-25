---
layout: "akamai"
page_title: "Akamai: Configuration"
subcategory: "Application Security"
description: |-
 Configuration
---



# akamai_appsec_configuration

**Scopes**: Security configuration

Returns information about all your security configurations, or returns information about a specific security configuration.

**Related API Endpoint**: [/appsec/v1/configs](https://techdocs.akamai.com/application-security/reference/get-configs)

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

data "akamai_appsec_configuration" "configurations" {
}

output "configuration_list" {
  value = data.akamai_appsec_configuration.configurations.output_text
}

data "akamai_appsec_configuration" "specific_configuration" {
  name = "Documentation"
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

This data source supports the following arguments:

- `name` (Optional). Name of the security configuration you want to return information for. If not included, information is returned for all your security configurations.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `config_id`. ID of the specified security configuration. Returned only if the `name` argument is included.
- `output_text`. Tabular report showing the `config_id,` `name`, `latest_version`, `version_active_in_staging`, and `version_active_in_production` values for all your security configurations.
- `latest_version`. Most-recent version number of the specified security configuration. Returned only if the `name` argument is included in the Terraform configuration file.
- `staging_version`. Version number of the specified security configuration currently active in staging. Returned only if the `name` argument is included in the Terraform configuration file.
- `production_version`. Version number of the specified security configuration currently active in production. Returned only if the `name` argument is included in the Terraform configuration file.