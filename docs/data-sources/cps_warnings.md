---
layout: akamai
subcategory: Certificate Provisioning System
---

# akamai_cps_warnings

Use the `akamai_cps_warnings` data source to return a map of all possible pre- and post-verification warnings. The map includes both the ID needed to acknowledge a warning and a brief description of the issue. 

CPS produces warnings during enrollment creation or after a client uploads the certificate. CPS won't process a change until you acknowledge all warnings.

You can use the warning IDs returned by this data source to acknowledge or auto-approve warnings. The `akamai_cps_third_party_enrollment` and `akamai_cps_upload_certificate` resources include arguments to help you do this.

## Basic usage

This example shows how to return a map of verification warnings:

```hcl

provider "akamai" {
  edgerc = "../config/edgerc"
  config_section = "shared_dns"
}

data "akamai_cps_warnings" "example" {}

```

## Argument reference

This data source supports does not support any arguments.


## Attributes reference

This data source returns this attribute:

  * `warnings` - Validation warnings for the current change you're making. Warnings display with an ID and a short description. Unless you auto-approve warnings, you need the ID to acknowledge the change. CPS won't process the change until you acknowledge these warnings. 

