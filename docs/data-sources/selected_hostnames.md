---
layout: "akamai"
page_title: "Akamai: SelectedHostnames"
subcategory: "APPSEC"
description: |-
 SelectedHostnames
---

# akamai_appsec_selected_hostnames

Use the `akamai_appsec_selected_hostnames` data source to retrieve a list of the hostnames that are currently protected under a given security configuration version.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

data "akamai_appsec_selected_hostnames" "selected_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
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

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `hostnames` - The list of selected hostnames.

* `hostnames_json` - The list of selected hostnames in json format.

* `output_text` - A tabular display of the selected hostnames.

