---
layout: "akamai"
page_title: "Akamai: Certificate Domain Validation"
subcategory: "CPS"
description: |-
  Domain Validation
---

# akamai_cps_dv_validation

Once you complete the Let’s Encrypt challenges, use the `akamai_cps_dv_validation` resource to send the acknowledgement to CPS and inform it that tokens are ready for validation.

## Example usage

Basic usage:

```hcl
resource "akamai_cps_dv_validation" "example" {
  enrollment_id = akamai_cps_dv_enrollment.dv.id
```
## Argument reference

The following arguments are supported:

* `enrollment_id` (Required) - Unique identifier for the DV certificate enrollment.

## Attributes reference

* `status` - The status of certificate validation.
