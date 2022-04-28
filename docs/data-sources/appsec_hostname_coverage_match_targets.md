---
layout: "akamai"
page_title: "Akamai: ApiHostnameCoverageMatchTargets"
subcategory: "Application Security"
description: |-
 ApiHostnameCoverageMatchTargets
---

# akamai_appsec_hostname_coverage_match_targets

**Scopes**: Hostname

Returns information about the API and website match targets used to protect a hostname. The returned information is described in the [Get the hostname coverage match targets](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getfailoverhostnames) section of the Application Security API.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/hostname-coverage/match-targets](https://techdocs.akamai.com/application-security/reference/get-coverage-match-targets)

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

data "akamai_appsec_hostname_coverage_match_targets" "match_targets" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostname  = "documentation.akamai.com"
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id`. (Required). Unique identifier of the security configuration associated with the hostname.
- `hostname` (Required). Name of the host you want to return information for. You can only return information for a single host and hostname at a time.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of the coverage information.
- `output_text`. Tabular report of the coverage information.