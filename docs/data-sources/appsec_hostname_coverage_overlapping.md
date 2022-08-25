---
layout: "akamai"
page_title: "Akamai: ApiHostnameCoverageOverlapping"
subcategory: "Application Security"
description: |-
 ApiHostnameCoverageOverlapping
---

# akamai_appsec_hostname_coverage_overlapping

**Scopes**: Security configuration; hostname

Returns information about any other configuration versions that contain a hostname found in the current configuration version. 

**Related API Endpoint**:[/appsec/v1/configs/{configId}/versions/{versionNumber}/hostname-coverage/overlapping](https://techdocs.akamai.com/application-security/reference/get-hostname-coverage-overlapping)

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

data "akamai_appsec_hostname_coverage_overlapping" "test" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  hostname  = "documentation.akamai.com"
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration you want to return information for.
- `hostname` (Required). Name of the host you want to return information for.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of the overlap information.
- `output_text`. Tabular report of the overlap information.