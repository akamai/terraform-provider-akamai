---
layout: "akamai"
page_title: "Akamai: MatchTarget"
subcategory: "APPSEC"
description: |-
  MatchTarget
---

# resource_akamai_appsec_match_target


The `resource_akamai_appsec_match_target` resource allows you to create or re-use MatchTargets.

If the MatchTarget already exists it will be used instead of creating a new one.

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


resource "akamai_appsec_match_target" "appsecmatchtarget" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    type =  "website"
    is_negative_path_match =  false
    is_negative_file_extension_match =  true
    default_file = "BASE_MATCH"
    hostnames =  ["example.com","www.example.net","m.example.com"]
    //file_paths =  ["/sssi/*","/cache/aaabbc*","/price_toy/*"]
    //file_extensions = ["wmls","jpeg","pws","carb","pdf","js","hdml","cct","swf","pct"]
    security_policy = "f1rQ_106946"
 
    bypass_network_lists = ["888518_ACDDCKERS","1304427_AAXXBBLIST"]
    
}
```

## Argument Reference

The following arguments are supported:
* `config_id`- (Required) The Configuration ID

* `version` - (Required) The Version Number of configuration

* `json` - (Optional) The JSON for configuration

* `type` - (Optional) The Version Number of configuration

* `sequence` - (Optional) The Version Number of configuration
* `is_negative_path_match` - (Optional) The Version Number of configuration
* `is_negative_file_extension_match` - (Optional) The Version Number of configuration
* `default_file` - (Optional) The Version Number of configuration
* `hostnames` - (Optional) The Version Number of configuration
* `file_paths` - (Optional) The Version Number of configuration
* `file_extensions` - (Optional) The Version Number of configuration
* `security_policy` - (Optional) The Version Number of configuration
* `bypass_network_lists` - (Optional) The Version Number of configuration

# Attributes Reference

The following are the return attributes:

*`targetid` - The TargetID

