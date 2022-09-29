---
layout: akamai
subcategory: EdgeWorkers
---

# akamai_edgeworker

Use the `akamai_edgeworker` data source to get an EdgeWorker for a given EdgeWorker ID.

## Example usage

This example returns the resource tier fields for the selected EdgeWorker ID:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_edgeworker" "test" {
  edgeworker_id = 3
  local_bundle  = "test_tmp/TestDataEdgeWorkersEdgeWorker/bundles/edgeworker_one_warning.tgz"
}
```

## Argument reference

The data source supports these arguments:

* `edgeworker_id` - (Required) The unique identifier of the EdgeWorker.
* `local_bundle` - (Optional) The path where the EdgeWorkers `.tgz` code bundle will be stored.

## Attributes reference

This data source returns these attributes:

* `name` - The EdgeWorker name.
* `group_id` - Defines the group association for the EdgeWorker.
* `resource_tier_id` - The unique identifier of a resource tier.
* `local_bundle_hash` - The local bundle hash for the EdgeWorker. It's used to identify content changes for the bundle.
* `version` - The bundle version.
* `warnings` - The list of warnings returned by EdgeWorker validation.
