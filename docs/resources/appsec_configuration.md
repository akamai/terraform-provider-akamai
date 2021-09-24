---
layout: "akamai"
page_title: "Akamai: Configuration"
subcategory: "Application Security"
description: |-
  Configuration
---

# akamai_appsec_configuration

**Scopes**: Contract and group

Creates a new WAP (Web Application Protector) or KSD (Kona Site Defender) security configuration. KSD security configurations start out empty (i.e., unconfigured), while WAP configurations are created using preset values. The contract referenced in the request body determines the type of configuration you can create.

In addition to manually creating a new configuration, you can use the `create_from_config_id` argument to clone an existing configuration.

**Related API Endpoint**: [/appsec/v1/configs](https://developer.akamai.com/api/cloud_security/application_security/v1.html#postconfigurations)

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

data "akamai_appsec_selectable_hostnames" "selectable_hostnames" {
  config_id = "Documentation"
}

resource "akamai_appsec_configuration" "create_config" {
  name        = "Documentation Test Configuration"
  description = "This configuration is used as a testing environment for the documentation team."
  contract_id = "5-2WA382"
  group_id    = 12198
  host_names  = ["documentation.akamai.com", "training.akamai.com"]
}

output "create_config_id" {
  value = akamai_appsec_configuration.create_config.config_id
}

// USE CASE: User wants to clone a new security configuration from an existing configuration and version.

resource "akamai_appsec_configuration" "clone_config" {
  name                  = "Documentation Test Configuration"
  description           = "This configuration is used as a testing environment for the documentation team."
  create_from_config_id = data.akamai_appsec_configuration.configuration.config_id
  create_from_version   = data.akamai_appsec_configuration.configuration.latest_version
  contract_id           = "5-2WA382"
  group_id              = 12198
  host_names            = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames
}

output "clone_config_id" {
  value = akamai_appsec_configuration.clone_config.config_id
}
```

## Argument Reference

This resource supports the following arguments:

- `name` (Required). Name of the new configuration.
- `description` (Required). Brief description of the new configuration.
- `create_from_config_id` (Optional). Unique identifier of the existing configuration being cloned in order to create the new configuration.
- `create_from_version` (Optional). Version number of the security configuration being cloned.
- `contract_id` (Required). Unique identifier of the Akamai contract t associated with the new configuration.
- `group_id` (Required). Unique identifier of the contract group associated with the new configuration.
- `host_names` (Required). JSON array containing the hostnames to be protected by the new configuration. You must specify at least one hostname in order to create a new configuration.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `config_id`. ID of the new security configuration.

