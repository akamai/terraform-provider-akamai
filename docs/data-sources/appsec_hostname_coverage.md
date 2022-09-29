---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_hostname_coverage

**Scopes**: Individual account

Returns information about the hostnames associated with your account. The returned data includes the hostname's protections, activation status, and other summary information. 

**Related API Endpoint**: [/appsec/v1/hostname-coverage](https://techdocs.akamai.com/application-security/reference/get-hostname-coverage)

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

// USE CASE: User wants to view hostname coverage data.

data "akamai_appsec_hostname_coverage" "hostname_coverage" {
}

output "hostname_coverage_list_json" {
  value = data.akamai_appsec_hostname_coverage.hostname_coverage.json
}

// USE CASE: User wants to display the returned data in a table.

output "hostname_coverage_list_output" {
  value = data.akamai_appsec_hostname_coverage.hostname_coverage.output_text
}
```

## Argument Reference

This data source does not support any arguments.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list of the hostname coverage information.
- `output_text`. Tabular report of the hostname coverage information.