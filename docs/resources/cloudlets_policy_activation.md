---
layout: akamai
subcategory: Cloudlets
---

# akamai_cloudlets_policy_activation

Use the `akamai_cloudlets_policy_activation` resource to activate a specific version of a Cloudlet policy. An activation deploys the version to either the Akamai staging or production network. You can activate a specific version multiple times if you need to.

Before activating on production, activate on staging first. This way you can detect any problems in staging before your changes progress to production.

## Example usage

Basic usage:

```hcl
resource "akamai_cloudlets_policy_activation" "example" {
  policy_id = 1234
  network = "staging"
  version = 1
  associated_properties = ["Property_1", "Property_2", "Property_3"]
}
```
If you're handling two `akamai_cloudlets_policy_activation` resources in the same configuration file with the same `policy_id`, but different `network` arguments (for example, `production` and `staging`), you need to add `depends_on` to the production resource. See the example:

```hcl
resource "akamai_cloudlets_policy_activation" "stag" {
  policy_id = 1234567
  network = "staging"
  version = 1
  associated_properties = ["Property_1","Property_2"]
}

resource "akamai_cloudlets_policy_activation" "prod" {
  policy_id = 1234567
  network = "production"
  version = 1
  associated_properties = ["Property_1","Property_2"]
  depends_on = [
    akamai_cloudlets_policy_activation.stag
  ]
}
```

## Argument reference

The following arguments are supported:

* `policy_id` - (Required) An identifier for the Cloudlet policy you want to activate.
* `network` - (Required) The network you want to activate the policy version on. For the Staging network, specify either `staging`, `stag`, or `s`. For the Production network, specify either `production`, `prod`, or `p`. All values are case insensitive.
* `version` - (Required) The Cloudlet policy version you want to activate.
* `associated_properties` - (Required) A set of property identifiers related to this Cloudlet policy. You can't activate a Cloudlet policy if it doesn't have any properties associated with it.

## Attribute reference

The following attributes are returned:

* `status` - The activation status for this Cloudlet policy.
