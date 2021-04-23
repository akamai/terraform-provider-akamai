---
layout: "akamai"
page_title: "Akamai: ApiHostnameCoverageOverlapping"
subcategory: "Application Security"
description: |-
 ApiHostnameCoverageOverlapping
---

# akamai_appsec_hostname_coverage_overlapping

Use the `akamai_appsec_hostname_coverage_overlapping` data source to retrieve information about the configuration versions that contain a hostname also included in the current configuration version. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#gethostnamecoverageoverlapping).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_hostname_coverage_overlapping" "test"  {
  config_id = 43253
  hostname = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `config_id`- (Required) The configuration ID.

* `hostname` - (Optional) The hostname for which to retrieve information.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list of the overlap information.

* `output_text` - A tabular display of the overlap information.

