---
layout: "akamai"
page_title: "Akamai: ConfigurationRename"
subcategory: "Application Security"
description: |-
  ConfigurationRename
---

# akamai_appsec_configuration_rename

The `akamai_appsec_configuration_rename` resource allows you to rename an existing security configuration.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to rename an existing configuration
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

resource "akamai_appsec_configuration_rename" "configuration" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  name = var.name
  description = var.description
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to be renamed.

* `name` - (Required) The new name to be given to the configuration.

* `description` - (Required) The description to be applied to the configuration.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* None

