---
layout: "akamai"
page_title: "Akamai: MatchTarget"
subcategory: "Application Security"
description: |-
  MatchTarget
---

# akamai_appsec_match_target

**Scopes**: Security configuration

Creates a match target associated with a security configuration. Match targets determine which security policy should apply to an API, hostname or path.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/match-targets](https://developer.akamai.com/api/cloud_security/application_security/v1.html#postmatchtargets)

## Example Usage

Basic usage:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

resource "akamai_appsec_match_target" "match_target" {
  config_id    = data.akamai_appsec_configuration.configuration.config_id
  match_target = file("${path.module}/match_targets.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the match target being modified.
- `match_target` (Required). Path to a JSON file containing one or more match target definitions. You can find a sample match target JSON file in the [Create a match target section](https://developer.akamai.com/api/cloud_security/application_security/v1.html#postmatchtargets) of the Application Security API documentation.

## Output Options

In addition to the arguments above, the following attribute is exported:

- `match_target_id`. ID of the match target.

