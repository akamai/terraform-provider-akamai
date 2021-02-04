---
layout: "akamai"
page_title: "Akamai: HostnameCoverage"
subcategory: "Application Security"
description: |-
 HostnameCoverage
---

# akamai_appsec_hostname_coverage

Use the `akamai_appsec_hostname_coverage` data source to retrieve a list of hostnames in the account with their current protections, activation statuses, and other summary information. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#8eb23096).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the hostname coverage data
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_hostname_coverage" "hostname_coverage" {
}

output "hostname_coverage_list_json" {
  value = data.akamai_appsec_hostname_coverage.hostname_coverage.json
}

//tabular data of hostname, status, hasMatchTarget
output "hostname_coverage_list_output" {
  value = data.akamai_appsec_hostname_coverage.hostname_coverage.output_text
}
```

## Argument Reference

The following arguments are supported:

* None

## Attributes Reference

The following attributes are exported:

* `json` - A JSON-formatted list of the hostname coverage information.

* `output_text` - A tabular display of the hostname coverage information.

