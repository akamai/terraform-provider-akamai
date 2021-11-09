---
layout: "akamai"
page_title: "Akamai: DataStream Activation History"
subcategory: "DataStream"
description: |-
 Activation History
---

# akamai_datastream_activation_history

Use the `akamai_datastream_activation_history` data source to list detailed information about the activation status changes for all versions of a stream.

## Example usage

This example returns the activation history for a provided stream ID:

```hcl
data "akamai_datastream_activation_history" ds {
  stream_id = 12345
}

output "ds_history_stream_id" {
  value = data.akamai_datastream_activation_history.ds.stream_id
}

output "ds_history_activations" {
  value = data.akamai_datastream_activation_history.ds.activations
}
```

## Argument reference

The data source supports this argument:

* `stream_id` - (Required) A stream's unique identifier.

## Attributes reference

This data source returns these attributes:

* `activations` - Detailed information about an activation status change for a version of a stream, including:
  * `created_by` - The user who activated or deactivated the stream.
  * `created_date` - The date and time of an activation status change.
  * `stream_id` - A stream's unique identifier.
  * `stream_version_id` - A stream version's unique identifier.
  * `is_active` -	Whether the version of the stream is active.
