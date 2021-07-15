---
layout: "akamai"
page_title: "Akamai: Certificate Domain Validation"
subcategory: "CPS"
description: |-
  Domain Validation
---

# akamai_cps_dv_validation

Once you complete the Let’s Encrypt challenges, optionally use the `akamai_cps_dv_validation` resource to send the acknowledgement to CPS and inform it that tokens are ready for validation. You can also wait for CPS to check for the tokens, which it does on a regular schedule. Next, CPS automatically deploys the certificate on Staging, and eventually on the Production network.

## Example usage

Basic usage:

```hcl
resource "akamai_cps_dv_validation" "example" {
  enrollment_id = akamai_cps_dv_enrollment.example.id
  sans = akamai_cps_dv_enrollment.example.sans
}
```
## Argument reference

The following arguments are supported:

* `enrollment_id` (Required) - Unique identifier for the DV certificate enrollment.
* `sans` - (Optional) The Subject Alternative Names (SAN) list for tracking changes on related enrollments. Whenever any SAN changes, the Akamai provider recreates this resource and sends another acknowledgement request to CPS.

## Attributes reference

* `status` - The status of certificate validation.
