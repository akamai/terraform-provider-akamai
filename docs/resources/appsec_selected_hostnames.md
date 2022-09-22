---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_selected_hostnames

**Scopes**: Security configuration

Modifies the list of hostnames protected under by a security configuration.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/selected-hostnames](https://techdocs.akamai.com/application-security/reference/put-selected-hostnames-per-config)

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

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_selected_hostnames" "appsecselectedhostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostnames = ["example.com"]
  mode      = "APPEND"
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the hostnames.
- `hostnames` (Required). JSON array of hostnames to be added or removed from the protected hosts list.
- `mode` (Required). Indicates how the `hostnames` array is to be applied. Allowed values are:
  - **APPEND**. Hosts listed in the `hostnames` array are added to the current list of selected hostnames.
  - **REPLACE**. Hosts listed in the `hostnames`  array overwrite the current list of selected hostnames: the “old” hostnames are replaced by the specified set of hostnames.
  - **REMOVE**, Hosts listed in the `hostnames` array are removed from the current list of select hostnames.