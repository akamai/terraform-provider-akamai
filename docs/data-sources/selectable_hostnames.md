---
layout: "akamai"
page_title: "Akamai: SelectableHostnames"
subcategory: "APPSEC"
description: |-
 SelectableHostnames
---

# akamai_appsec_selectable_hostnames

Use the `akamai_appsec_selectable_hostnames` data source to retrieve the list of hostnames that may be protected under a given security configuration version.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

data "akamai_appsec_selectable_hostnames" "selectable_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
}

output "selectable_hostnames" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames
}

output "selectable_hostnames_json" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.hostnames_json
}

output "selectable_hostnames_output_text" {
  value = data.akamai_appsec_selectable_hostnames.selectable_hostnames.output_text
}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `hostnames` - The list of selectable hostnames.

* `hostnames_json` - The list of selectable hostnames in json format.

* `output_text` - A tabular display of the selectable hostnames showing the name and config_id of the security configuration under which the host is protected in production, or '-' if the host is not protected in production.

