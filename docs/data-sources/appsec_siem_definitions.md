---
layout: "akamai"
page_title: "Akamai: SiemDefinitions"
subcategory: "Application Security"
description: |-
 SiemDefinitions
---

# akamai_appsec_siem_definitions

The `akamai_appsec_siem_definitions` data source allows you to retrieve information about the available SIEM versions, or about a specific SIEM version. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getsiemversions).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the siem settings with a given security configuration
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_siem_definitions" "siem_definitions" {
}

output "siem_definitions_json" {
  value = data.akamai_appsec_siem_definitions.siem_definitions.json
}

output "siem_definitions_output" {
  value = data.akamai_appsec_siem_definitions.siem_definitions.output_text
}

data "akamai_appsec_siem_definitions" "siem_definition" {
  siem_definition_name = var.siem_definition_name
}

output "siem_definition_id" {
  value = data.akamai_appsec_siem_definitions.siem_definition.id
}
```

## Argument Reference

The following arguments are supported:

* `siem_definition_name`- (Optional) The name of a specific SIEM definition for which to retrieve information.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list of the SIEM version information.

* `output_text` - A tabular display showing the ID and name of each SIEM version.

