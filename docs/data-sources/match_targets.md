---
layout: "akamai"
page_title: "Akamai: MatchTargets"
subcategory: "APPSEC"
description: |-
 MatchTargets
---

# akamai_appsec_match_targets

Use `akamai_appsec_match_targets` data source to retrieve a match_targets id.

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


resource "akamai_appsec_match_targets" "appsecmatchtargets" {
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


# Attributes Reference

The following are the return attributes:

*`output_text` - The match targets in formatted text

