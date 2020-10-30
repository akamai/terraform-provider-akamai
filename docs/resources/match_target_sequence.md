---
layout: "akamai"
page_title: "Akamai: MatchTargetSequence"
subcategory: "APPSEC"
description: |-
  MatchTargetSequence
---

# resource_akamai_appsec_match_target_sequence


The `resource_akamai_appsec_match_target_sequence` resource allows you to specify the order in which match targets are applied within a given security configuration and version.


## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

resource "akamai_appsec_match_target" "match_target_sequence" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  type =  "website"
  // json =  file("${path.module}/match_targets.json")
  sequence_map = {
	  2971336 = 1
	  2052813 = 2
  }  
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `json` - The name of a JSON file containing the sequence of all match targets defined for the specified security configuration and version. If not specified, the match target sequence must be specified using the `type` and `sequence_map` parameters described below.

* `type` - (Required) Describes the type of match target, either website or api. Must not be specified if `json` is specified. Must not be specified if `json` is specified.

* `sequence_map` - (Required) A list specifying the IDs and sequence numbers of all match targets defined within the specified security configuration and version. Must not be specified if `json` is specified.


In addition to the arguments above, the following attributes are exported:

* None




