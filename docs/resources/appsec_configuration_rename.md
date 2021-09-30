---
layout: "akamai"
page_title: "Akamai: ConfigurationRename"
subcategory: "Application Security"
description: |-
  ConfigurationRename
---

# akamai_appsec_configuration_rename

**Scopes**: Security configuration

Renames an existing security configuration. 
Note that you can only change the configuration name.
The ID assigned to a security configuration can not be modified.

**Related API Endpoint**: [/appsec/v1/configs/{configId}](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putrenameconfiguration)

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

// USE CASE: User wants to rename an existing security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_configuration_rename" "configuration" {
  config_id   = data.akamai_appsec_configuration.configuration.config_id
  name        = "Documentation and Training Configuration"
  description = "This configuration is by both the documentation team and the training team."
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configurating being renamed.
- `name` (Required). New name for the security configuration.
- `description` (Required). Brief description of the security configuration.

