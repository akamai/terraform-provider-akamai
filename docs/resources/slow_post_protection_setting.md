---
layout: "akamai"
page_title: "Akamai: SlowPostProtectionSetting"
subcategory: "APPSEC"
description: |-
  SlowPostProtectionSetting
---

# resource_akamai_appsec_slow_post_protection_setting


The `resource_akamai_appsec_slow_post_protection_setting` resource allows you to create or re-use SlowPostProtectionSettings.

If the SlowPostProtectionSetting already exists it will be used instead of creating a new one.

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


resource "akamai_appsec_slow_post_protection_settings" "appsecslowpostprotectionsettings" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
    slow_rate_action = "alert"                        
    slow_rate_threshold_rate = 10
    slow_rate_threshold_period = 30
    duration_threshold_timeout = 20
}
```

## Argument Reference

The following arguments are supported:
* `config_id`- (Required) The Configuration ID

* `versionnumber` - (Required) The Version Number of configuration

# Attributes Reference

The following are the return attributes:

*`targetid` - The TargetID

