---
layout: "akamai"
page_title: "Akamai: MatchTargetSequence"
subcategory: "Application Security"
description: |-
  MatchTargetSequence
---

# akamai_appsec_match_target_sequence

**Scopes**: Security configuration

Specifies the order in which match targets are applied within a security configuration. As a general rule, you should process broader and more-general match targets first, gradually working your way down to more granular and highly-specific targets.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/match-targets/sequence](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putsequence)

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

resource "akamai_appsec_match_target_sequence" "match_target_sequence" {
  config_id             = data.akamai_appsec_configuration.configuration.config_id
  match_target_sequence = file("${path.module}/match_targets_sequence.json")
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the match target sequence being modified.
- `match_target_sequence` (Required). Path to a JSON file containing the processing sequence for all the match targets defined for the security configuration. You can find a sample match target sequence JSON file in the [Modify match target order](https://developer.akamai.com/api/cloud_security/application_security/v1.html#matchtargetorder) section of the Application Security API documentation.

