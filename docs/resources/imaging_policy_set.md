---
layout: akamai
subcategory: Image and Video Manager
---

# akamai_imaging_policy_video

Use the `akamai_imaging_policy_set` resource to define a policy set.

## Basic usage

This example returns the policy set details:

```hcl
resource "akamai_imaging_policy_set" "example_policy_set" {
    contract_id            = "1234"
    name                   = "image_policyset"
    region                 = "US"
    type                   = "IMAGE"
}
```

## Argument reference

This resource supports these arguments:
* `contract_id` - (Required) The unique identifier for the Akamai Contract containing the policy set.
* `name` - (Required) A friendly name for the policy set.
* `region` - (Required) The geographic region for which the media using this policy set is optimized: `US`, `EMEA`, `ASIA`, `AUSTRALIA`, `JAPAN` or `CHINA`
* `type` - (Required) The type of media managed by this policy set: `IMAGE` or `VIDEO`
