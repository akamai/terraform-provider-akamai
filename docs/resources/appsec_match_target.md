---
layout: "akamai"
page_title: "Akamai: MatchTarget"
subcategory: "Application Security"
description: |-
  MatchTarget
---

# akamai_appsec_match_target


The `akamai_appsec_match_target` resource allows you to create or modify a match target associated with a given security configuration.


## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

resource "akamai_appsec_match_target" "match_target" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  match_target =  file("${path.module}/match_targets.json")
}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `match_target` - (Required) The name of a JSON file containing one or more match target definitions ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#postmatchtargets)).

## Attribute Reference

In addition to the arguments above, the following attribute is exported:

* `match_target_id` - The ID of the match target.

