---
layout: "akamai"
page_title: "Akamai: Configuration"
subcategory: "APPSEC"
description: |-
 Configuration
---

# akamai_appsec_configuration

Use `akamai_appsec_configuration` data source to retrieve a configuration id.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}


data "akamai_appsec_configuration" "appsecconfiguration" {
    name = "Akamai Tools"
   }

output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}

output "configsedgelatestversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}

output "configsedgeconfiglist" {
  value = data.akamai_appsec_configuration.appsecconfigedge.output_text
}

output "configsedgeconfigversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.version
}

```

## Argument Reference

The following arguments are supported:

* `name`- (Optional) The Configuration Name

* `version` - (Optional) The Version Number of configuration

# Attributes Reference

The following are the return attributes:

* `config_id` - Configuration data

* `output_text` - Configuration list in tabular format

