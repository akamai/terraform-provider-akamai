---
layout: "akamai"
page_title: "Akamai: MatchTargets"
subcategory: "Application Security"
description: |-
 MatchTargets
---

# akamai_appsec_match_targets

Use the `akamai_appsec_match_targets` data source to retrieve information about the match targets associated with a given configuration version, or about a specific match target.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view the match targets associated with a given security configuration
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_match_targets" "match_targets" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
}
output "match_targets" {
  value = data.akamai_appsec_match_targets.match_targets.output_text
}

// USE CASE: user wants to see a single match target associated with a given security configuration version
data "akamai_appsec_match_targets" "match_target" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  match_target_id = var.match_target_id
}
output "match_target_output" {
  value = data.akamai_appsec_match_targets.match_target.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `match_target_id` - (Optional) The ID of the match target to use. If not supplied, information about all match targets is returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the ID and Policy ID of all match targets associated with the specified security configuration and version, or of the specific match target if `match_target_id` was supplied.

