---
layout: "akamai"
page_title: "Akamai: ConfigurationClone"
subcategory: "Application Security"
description: |-
  ConfigurationClone
---

# resource_akamai_appsec_configuration_clone

The `resource_akamai_appsec_configuration_clone` resource allows you to create a new version of a given security configuration.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to clone a new configuration from an existing one
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

// USE CASE: user wants to see contract group details in an account
data "akamai_appsec_contracts_groups" "contracts_groups" {
  contractid = var.contractid
  groupid = var.groupid
}

data "akamai_appsec_selectable_hostnames" "selectable_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
}

resource "akamai_appsec_configuration_clone" "clone_config" {
  create_from_config_id = data.akamai_appsec_configuration.configuration.config_id
  create_from_version = data.akamai_appsec_configuration.configuration.latest_version
  name = var.name
  description = var.description
  contract_id = data.akamai_appsec_contracts_groups.contracts_groups.default_contractid
  group_id = data.akamai_appsec_contracts_groups.contracts_groups.default_groupid
  host_names = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames
}

output "clone_config_id" {
  value = akamai_appsec_configuration_clone.clone_config.config_id
}

output "clone_config_version" {
  value = akamai_appsec_configuration_clone.clone_config.version
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to be applied to the new configuration.

* `description` - (Required) A description of the new configuration.

* `create_from_config_id` - (Required) The ID of the configuration to be cloned.

* `create_from_version` - (Required) The version number of the configuration to be cloned.

* `contract_id` - (Required)  The contract id to use.

* `group_id` - (Required)  The group id to use.

* `host_names` - (Required)  The hostnames to be protected under the new configuration.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `config_id` - The ID of the newly created configuration.

* `version` - The version number of the newly created configuration.

