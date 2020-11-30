---
layout: "akamai"
page_title: "Akamai: MatchTargets"
subcategory: "Application Security"
description: |-
 MatchTargets
---

# akamai_appsec_match_targets

Use the `akamai_appsec_match_targets` data source to retrieve information about the match targets associated with a given configuration version.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

data "akamai_appsec_match_targets" "match_targets" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
}

output "match_targets" {
  value = data.akamai_appsec_match_targets.match_targets.output_text
}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the ID and Policy ID of all match targets associated with the specified security configuraton and version.

