---
layout: "akamai"
page_title: "Akamai: MatchTarget"
subcategory: "APPSEC"
description: |-
  MatchTarget
---

# resource_akamai_appsec_match_target


The `resource_akamai_appsec_match_target` resource allows you to create or modify a match target associated with a given security configuration and version.


## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

resource "akamai_appsec_match_target" "match_target_1" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  json =  file("${path.module}/match_targets.json")
}

resource "akamai_appsec_match_target" "match_target_2" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  type =  "website"
  is_negative_path_match =  false
  is_negative_file_extension_match =  true
  default_file = "NO_MATCH"
  hostnames =  ["example.com","www.example.net","n.example.com"]
  file_paths =  ["/sssi/*","/cache/aaabbc*","/price_toy/*"]
  file_extensions = ["wmls","jpeg","pws","carb","pdf","js","hdml","cct","swf","pct"]
  security_policy = "crAP_75829"
  bypass_network_lists = ["12345_FOO","67890_BAR"]
}


```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `json` - The name of a JSON file containing one or more match target definitions. If not specified, the match target must be specified using the additional parameters listed below.

  * `type`
  * `is_negative_path_match`
  * `is_negative_file_extension_match`
  * `default_file`
  * `hostnames`
  * `file_paths`
  * `file_extensions`
  * `security_policy`
  * `bypass_network_lists`

* `type` - (Required) Describes the type of match target, either website or api. Must not be specified if `json` is specified.

* `is_negative_path_match` - Describes whether the match target applies when a match is found in the specified paths or when a match isn’t found. Must not be specified if `json` is specified.

* `is_negative_file_extension_match` - Describes whether the match target applies when a match is found in the specified fileExtensions or when a match isn’t found. Must not be specified if `json` is specified.

* `default_file` - Describes the rule to match on paths. Either NO_MATCH to not match on the default file, BASE_MATCH to match only requests for top-level hostnames ending in a trailing slash, or RECURSIVE_MATCH to match all requests for paths that end in a trailing slash. Must not be specified if `json` is specified.

* `hostnames` - The hostnames to match the request on. Must not be specified if `json` is specified.

* `file_paths` - The path used in the path match. Must not be specified if `json` is specified.

* `file_extensions` - The file extensions used in the path match. Must not be specified if `json` is specified.

* `security_policy` - (Required) The security policy associated with the match target. Must not be specified if `json` is specified.

* `bypass_network_lists` - The list of network list identifiers and names. Must not be specified if `json` is specified.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `target_id` - The ID of the match target.



