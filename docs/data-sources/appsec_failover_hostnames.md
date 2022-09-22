---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_failover_hostnames

**Scopes**: Security configuration

Returns a list of the failover hostnames in a configuration. 

**Related API Endpoint**: [/appsec/v1/configs/{configId}/failover-hostnames](https://techdocs.akamai.com/application-security/reference/get-failover-hostnames)

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

// USE CASE: User wants to view the failover hostnames for a security configuration.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_failover_hostnames" "failover_hostnames" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "failover_hostnames" {
  value = data.akamai_appsec_failover_hostnames.failover_hostnames.hostnames
}

output "failover_hostnames_output" {
  value = data.akamai_appsec_failover_hostnames.failover_hostnames.output_text
}

output "failover_hostnames_json" {
  value = data.akamai_appsec_failover_hostnames.failover_hostnames.json
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the failover hosts.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `hostnames`. List of the failover hostnames.
- `json`. JSON-formatted list of the failover hostnames.