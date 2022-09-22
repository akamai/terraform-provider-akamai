---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_match_targets

**Scopes**: Security configuration; match target

Returns information about your match targets. Match targets determine which security policy should apply to an API, hostname, or path.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/match-targets{?policyId,includeChildObjectName}](https://techdocs.akamai.com/application-security/reference/get-match-targets)

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

// USE CASE: User wants to view the match targets associated with a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_match_targets" "match_targets" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}
output "match_targets" {
  value = data.akamai_appsec_match_targets.match_targets.output_text
}

// USE CASE: User wants to view a single match target.

data "akamai_appsec_match_targets" "match_target" {
  config_id       = data.akamai_appsec_configuration.configuration.config_id
  match_target_id = "2712938"
}
output "match_target_output" {
  value = data.akamai_appsec_match_targets.match_target.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the match targets.
- `match_target_id` (Optional). Unique identifier of the match target you want to return information for. If not included, information is returned for all your match targets.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the ID and security policy ID of your match targets.
- `json`. JSON-formatted list of the match target information.