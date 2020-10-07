---
layout: "akamai"
page_title: "Akamai: ConfigurationClone"
subcategory: "APPSEC"
description: |-
  ConfigurationClone
---

# resource_akamai_appsec_configuration_clone


The `resource_akamai_appsec_configuration_clone` resource allows you to create or re-use ConfigurationClones.

If the ConfigurationClone already exists it will be used instead of creating a new one.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}
 data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Example for EDGE"
  
}



output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}


resource "akamai_appsec_configuration_clone" "appsecconfigurationclone" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    create_from_version = data.akamai_appsec_configuration.appsecconfigedge.latest_version 
    rule_update  = false
   }

```

## Argument Reference

The following arguments are supported:
* `config_id`- (Required) The Configuration ID

* `create_from_version` - (Required) The Version Number of configuration

* `rule_update` - (Optional) Update Rules Flag

# Attributes Reference

The following are the return attributes:

* `configcloneid` - Id of cloned configuration

* `version` - Version of cloned configuration

