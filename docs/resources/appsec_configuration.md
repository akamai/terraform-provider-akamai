---
layout: "akamai"
page_title: "Akamai: Configuration"
subcategory: "Application Security"
description: |-
  Configuration
---

# resource_akamai_appsec_configuration

The `resource_akamai_appsec_configuration` resource allows you to create a new WAP or KSD security configuration. KSD security configurations start out empty, and WAP configurations are created with preset values. The contract you pass in the request body determines which product you use. You can edit the default settings included in the WAP configuration, but you’ll need to run additional operations in this API to select specific protections for KSD. Your KSD configuration needs match targets and protection settings before it can be activated. 

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to create a new config
data "akamai_appsec_contract_groups" "contract_groups" {
}

data "akamai_appsec_selectable_hostnames" "selectable_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
}

resource "akamai_appsec_configuration" "create_config" {
  name = var.name
  description = var.description
  contract_id= data.akamai_appsec_contract_groups.contract_groups.default.contract_id
  group_id  = data.akamai_appsec_contract_groups.contract_groups.default.group_id
  host_names = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames
}

output "create_config_id" {
  value = akamai_appsec_configuration.create_config.config_id
}

output "create_config_version" {
  value = akamai_appsec_configuration.create_config.version
}
```

## Argument Reference

The following arguments are supported:

* `name`- (Required) The name to be assigned to the configuration.

* `description` - (Required) A description of the configuration.

* `contract_id` - (Required) The contract ID of the configuration.

* `group_id` - (Required) The group ID of the configuration>

* `hostnames` - (Required) The list of hostnames protected by this security configuration.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `config_id` - (Required) The ID of the security configuration.

* `version` - (Required) The latest version of the security configuration.

