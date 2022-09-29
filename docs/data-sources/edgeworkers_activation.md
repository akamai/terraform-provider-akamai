---
layout: akamai
subcategory: EdgeWorkers
---

# akamai_edgeworker_activation

Use the `akamai_edgeworker_activation` data source to fetch the latest activation for a given EdgeWorker ID.

## Example usage

This example returns the latest activation on the staging network:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_edgeworker_activation" "test" {
  edgeworker_id = 1
  network       = "STAGING"
}
```

## Argument reference

The data source supports these arguments:

* `edgeworker_id` - (Required) The unique identifier of the EdgeWorker.
* `network` - (Required) The network from where the activation information will be fetched.

## Attributes reference

This data source returns these attributes:

* `activation_id` - The unique identifier of the activation.
* `version` - The EdgeWorker version of the latest activation.
