---
layout: "akamai"
page_title: "Akamai: ApiHostnameCoverageMatchTargets"
subcategory: "Application Security"
description: |-
 ApiHostnameCoverageMatchTargets
---

# akamai_appsec_hostname_coverage_match_targets

Use the `akamai_appsec_hostname_coverage_match_targets` data source to retrieve information about the API and website match targets that protect a hostname. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#gethostnamecoveragematchtargets).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_hostname_coverage_match_targets" "match_targets" {
  config_id = 43253
  hostname = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `config_id`- (Required) The configuration ID.

* `hostname` - (Required) The hostname for which to retrieve information.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list of the coverage information.

* `output_text` - A tabular display of the coverage information.

