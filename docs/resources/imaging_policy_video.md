---
layout: akamai
subcategory: Image and Video Manager
---

# akamai_imaging_policy_video (Beta)

Use the `akamai_imaging_policy_video` resource to list details about a policy.

## Basic usage

This example returns the policy details based on the policy ID and optionally, a version:

```hcl
resource "akamai_imaging_policy_video" "example" {
    activate_on_production = false
    version                = 1
    contract_id            = "1234"
    policy_id              = "imgpolicy1234"
    policyset_id           = "akamai_imaging_policy_set.policy_set_name.id"
    json                   = file("policy.json")  
}
```

## Argument reference

This resource supports these arguments:
* `activate_on_production` - (Optional) With this flag set to `false`, the user can perform modifications on staging without affecting the version already saved to production.
With this flag set to `true`, the policy will also be saved on the production network.
It is possible to change it back to `false` only when there are any changes to the policy qualifying it for the new version.
It should be set to false whenever there are changes to policy to ensure that the change is deployed to and tested on staging first.
* `contract_id` - (Required) The nique identifier for the Akamai Contract containing the policy set.
* `policy_id` - (Required) The unique identifier of a policy.
It is not possible to modify the id of the policy.
* `policyset_id` - (Required) The unique identifier for the Image & Video Manager policy set.
* `json` - (Required) A JSON encoded policy.


## Attributes reference

This resource returns this attribute:

* `version` - The version number of the policy.
