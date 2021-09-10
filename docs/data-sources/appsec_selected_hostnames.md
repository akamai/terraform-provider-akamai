---
layout: "akamai"
page_title: "Akamai: SelectedHostnames"
subcategory: "Application Security"
description: |-
 SelectedHostnames
---

# akamai_appsec_selected_hostnames

**Scopes**: Security configuration

Returns a list of the hostnames currently protected by the specified security configuration.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/selected-hostnames](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getselectedhostnames)

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

data "akamai_appsec_selected_hostnames" "selected_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "selected_hostnames" {
  value = data.akamai_appsec_selected_hostnames.selected_hostnames.hostnames
}

output "selected_hostnames_json" {
  value = data.akamai_appsec_selected_hostnames.selected_hostnames.hostnames_json
}

output "selected_hostnames_output_text" {
  value = data.akamai_appsec_selected_hostnames.selected_hostnames.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the protected hosts.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `hostnames`. List of selected hostnames.
- `hostnames_json`. JSON-formatted list of selected hostnames.
- `output_text`. Tabular report of the selected hostnames.

