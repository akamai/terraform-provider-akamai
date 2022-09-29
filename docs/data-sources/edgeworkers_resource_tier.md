---
layout: akamai
subcategory: EdgeWorkers
---

# akamai_edgeworkers_resource_tier

Use the `akamai_edgeworkers_resource_tier` data source to list the available resource tiers for a specific contract ID. The resource tier defines the resource consumption [limits](https://techdocs.akamai.com/edgeworkers/docs/resource-tier-limitations) for an EdgeWorker ID.

## Example usage

This example returns the resource tier fields for an EdgeWorker ID:

```hcl
data "akamai_edgeworkers_resource_tier" "example" {
  contract_id        = "1-ABC"
  resource_tier_name = "Basic Compute"
}
```

## Argument reference

The data source supports this argument:

* `contract_id` - (Required) Unique identifier of a contract.
* `resource_tier_name` - (Required) Unique name of the resource tier.

## Attributes reference

This data source returns these attributes:

* `resource_tier_id` - Unique identifier of the resource tier.
