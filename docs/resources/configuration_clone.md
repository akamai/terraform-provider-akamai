---
layout: "akamai"
page_title: "Akamai: ConfigurationClone"
subcategory: "APPSEC"
description: |-
  ConfigurationClone
---

# resource_akamai_appsec_configuration_clone


The `resource_akamai_appsec_configuration_clone` resource allows you to create a new version of a security configuration by cloning an existing version.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

resource "akamai_appsec_configuration_version_clone" "clone" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  create_from_version = data.akamai_appsec_configuration.configuration.latest_version
  rule_update  = false
}

output "clone_version" {
  value = akamai_appsec_configuration_version_clone.clone.version
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `create_from_version` - (Required) The version number of the security configuration to clone.

* `rule_update` - A boolean indicating whether to update the rules of the new version. If not supplied, False is assumed.

## Attribute Reference

In addition to the arguments above, the following attribute is exported:

* `version` - The number of the cloned version.

