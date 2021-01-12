---
layout: "akamai"
page_title: "Akamai: FailoverHostnames"
subcategory: "Application Security"
description: |-
 FailoverHostnames
---

# akamai_appsec_failover_hostnames

The `akamai_appsec_failover_hostnames` data source allows you to retrieve a list of the failover hostnames in a configuration. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getfailoverhostnames).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the failover hostnames in a given security configuration
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
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

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

## Attributes Reference

In addition to the argument above, the following attributes are exported:

* `hostnames` - A list of the failover hostnames.

* `json` - A JSON-formatted list of the failover hostnames.

* `output_text` - A tabular display showing the failover hostnames.

