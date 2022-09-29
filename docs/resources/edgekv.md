---
layout: akamai
subcategory: EdgeWorkers
---

# akamai_edgekv

The `akamai_edgekv` resource lets you control EdgeKV database functions outside EdgeWorkers JavaScript code. Refer to the [EdgeKV documentation](https://techdocs.akamai.com/edgekv/docs/welcome-to-edgekv) for more information.

## Example usage

Basic usage:

```hcl
resource "akamai_edgekv" "test_staging" {
  network              = "staging"
  namespace_name       = "Marketing"
  retention_in_seconds = 15724800
  group_id             = 4284
  geo_location         = "US"
  initial_data {
    key   = "lang"
    value = "English"
    group = "translations"
  }
}
```

## Argument reference

This resource supports these arguments:

* `namespace_name` - (Required) The name of the namespace.
* `network` - (Required) The network you want to activate the EdgeKV database on. For the Staging network, specify either `STAGING`, `STAG`, or `S`. For the Production network, specify either `PRODUCTION`, `PROD`, or `P`. All values are case insensitive.
* `group_id` - (Required) The `group ID` for the EdgeKV namespace. This numeric value will be required in the next EdgeKV API version.
* `retention_in_seconds` - (Required) Retention period for data in this namespace, or 0 for indefinite. An update of this value will just affect new EdgeKV items.
* `geo_location` - (Optional) Storage location for data when creating a namespace on the production network. This can help optimize performance by storing data where most or all of your users are located. The value defaults to `US` on the `STAGING` and `PRODUCTION` networks. For a list of supported geoLocations on the `PRODUCTION` network refer to the [EdgeKV documentation](https://techdocs.akamai.com/edgekv/docs/edgekv-data-model#namespace).
* `initial_data` - (Optional) List of key-value pairs called items to initialize the namespace. These items are valid only for database creation, updates are ignored.

## Attributes reference

There are no supported arguments for this resource.
