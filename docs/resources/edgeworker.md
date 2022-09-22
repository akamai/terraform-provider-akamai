---
layout: akamai
subcategory: EdgeWorkers
---

# akamai_edgeworker

The `akamai_edgeworker` resource lets you deploy custom code on thousands of edge servers and apply logic that creates powerful web experiences.

## Example usage

Basic usage:

```hcl
resource "akamai_edgeworker" "ew" {
  group_id         = 72297
  name             = "Ew_42"
  resource_tier_id = 100
  local_bundle     = var.bundle_path
}
```

## Argument reference

This resource supports these arguments:

* `name` - (Required) The name of the EdgeWorker ID.
* `group_id` - (Required) Identifies a group to assign to the EdgeWorker ID.
* `resource_tier_id` - (Required) Unique identifier of the resource tier.
* `local_bundle` - (Optional) The path to the EdgeWorkers code bundle.

## Attributes reference

* `edgeworker_id` - Unique identifier for an EdgeWorker ID.
* `local_bundle_hash` - A SHA-256 hash digest of the EdgeWorkers code bundle.
* `version` - Unique identifier for a specific EdgeWorker version.
* `warnings` - List of validation warnings.
