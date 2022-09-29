---
layout: akamai
subcategory: EdgeWorkers
---

# akamai_edgeworkers_activation

Use the `akamai_edgeworkers_activation` resource to activate a specific EdgeWorker version. An activation deploys the version to either the Akamai staging or production network.

Before activating on production, activate on staging first. This way you can detect any problems in staging before your changes progress to production.

## Example usage

Basic usage:

```hcl
resource "akamai_edgeworkers_activation" "test" {
  edgeworker_id = 1234
  network       = "STAGING"
  version       = "test1"
}
```

## Argument reference

The following arguments are supported:

* `edgeworker_id` - (Required) A unique identifier for the EdgeWorker ID you want to activate.
* `version` - (Required) The EdgeWorker version you want to activate.
* `network` - (Required) The network you want to activate the policy version on. For the Staging network, specify either `STAGING`, `STAG`, or `S`. For the Production network, specify either `PRODUCTION`, `PROD`, or `P`. All values are case insensitive.

-> **Note** You can use the staging network to validate the behavior of your EdgeWorkers code bundle. Once you've tested the functionality, you can activate it on the production network.

## Attribute reference

The following attributes are returned:

* `activation_id` - (Required) Unique identifier of the activation.
