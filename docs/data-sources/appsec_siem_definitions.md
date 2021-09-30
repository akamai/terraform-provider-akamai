---
layout: "akamai"
page_title: "Akamai: SiemDefinitions"
subcategory: "Application Security"
description: |-
 SiemDefinitions
---

# akamai_appsec_siem_definitions

**Scopes**: SIEM definition

Returns information about your SIEM (Security Information and Event Management) versions. The returned information is described in the [Get SIEM versions](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsiemversions) section of the Application Security API.

**Related API Endpoint**: [/appsec/v1/siem-definitions](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsiemversions)

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

// USE CASE: User wants to view the SIEM settings for a security configuration.

data "akamai_appsec_siem_definitions" "siem_definitions" {
}

output "siem_definitions_json" {
  value = data.akamai_appsec_siem_definitions.siem_definitions.json
}

output "siem_definitions_output" {
  value = data.akamai_appsec_siem_definitions.siem_definitions.output_text
}

data "akamai_appsec_siem_definitions" "siem_definition" {
  siem_definition_name = "SIEM Version 01"
}

output "siem_definition_id" {
  value = data.akamai_appsec_siem_definitions.siem_definition.id
}
```

## Argument Reference

This data source supports the following arguments:

- `siem_definition_name` (Optional). Name of the SIEM definition you want to return information for. If not included, information is returned for all your SIEM definitions.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of the SIEM version information.
- `output_text`. Tabular report showing the ID and name of each SIEM version.

