---
layout: "akamai"
page_title: "Akamai: SlowPostProtectionSettings"
subcategory: "APPSEC"
description: |-
 SlowPostProtectionSettings
---

# akamai_appsec_slow_post_protection_settings

Use `akamai_appsec_slow_post_protection_settings` data source to retrieve a slow_post_protection_settings id.

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


data "akamai_appsec_slow_post" "appsecslowpostprotectionsettings" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
}

output "configsedge_post_output_text" {
  value = data.akamai_appsec_slow_post.appsecslowpostprotectionsettings.output_text
}

```

## Argument Reference

The following arguments are supported:

* `config_id`- (Required) The Configuration ID

* `version` - (Required) The Version Number of configuration

# Attributes Reference

The following are the return attributes:

*`output_text` - The slow post settings in  formatted text

