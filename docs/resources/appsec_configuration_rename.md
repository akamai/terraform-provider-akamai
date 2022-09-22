---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_configuration_rename

**Scopes**: Security configuration

Renames an existing security configuration.
Note that you can change only the configuration name. You can't modify the ID assigned to a security configuration.

**Related API Endpoint**: [/appsec/v1/configs/{configId}](https://techdocs.akamai.com/application-security/reference/put-config)

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