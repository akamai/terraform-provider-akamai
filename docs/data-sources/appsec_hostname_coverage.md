---
layout: "akamai"
page_title: "Akamai: HostnameCoverage"
subcategory: "Application Security"
description: |-
 ApiHostnameCoverage
---

# akamai_appsec_hostname_coverage

The `akamai_appsec_hostname_coverage` data source allows you to retrieve a list of hostnames in the account with their current protections, activation statuses, and other summary information. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#8eb23096).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_hostname_coverage" "hostname_coverage" {
}
```

## Argument Reference
The following arguments are supported:

* None

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list of the hostname coverage information.

* `output_text` - A tabular display of the hostname coverage information.

