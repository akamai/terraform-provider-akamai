---
layout: akamai
subcategory: DataStream
---

# akamai_datastreams

Use the `akamai_datastreams` data source to list details about the DataStream configuration.

## Example usage

This example returns stream data for a specific group ID: 
```hcl
locals {
  concrete_stream = [for stream in data.akamai_datastreams.stream_list_in_group.streams : stream if stream.stream_name == "Concrete stream name"][0]
}

data "akamai_datastreams" "stream_list_in_group" {
  group_id = "1234"
}
```

## Argument reference

The data source supports this argument:

* `group_id` - (Optional) Unique identifier of the group that can access the product.

## Attributes reference

This data source returns these attributes:

* `streams` - Returns the latest versions of the stream configurations for all groups within in your account. You can use the `group_id` parameter to view the latest versions of all configurations in a specific group. 
  * `activation_status` - The activation status of the stream. These are possible values: `ACTIVATED`, `DEACTIVATED`, `ACTIVATING`, `DEACTIVATING`, or `INACTIVE`. See the [Activate a stream](https://techdocs.akamai.com/datastream2/reference/put-stream-activate) and [Deactivate a stream](https://techdocs.akamai.com/datastream2/reference/put-stream-deactivate) operations.
  * `archived` - Whether the stream is archived.
  * `connectors` - The connector where the stream sends logs. 
  * `contract_id` - Identifies the contract that the stream is associated with.
  * `created_by` - The user who created the stream.
  * `created_date` - The date and time when the stream was created in this format: `14-07-2020 07:07:40 GMT`.
  * `current_version_id` - Identifies the current version of the stream.
  * `errors` - Objects that may indicate stream failure errors. Learn more about [Errors](https://techdocs.akamai.com/datastream2/reference/errors).
    * `detail` - A message informing about the status of the failed stream.
    * `title` - A descriptive label for the type of error.
    * `type` - Identifies the error type, either `ACTIVATION_ERROR` or `UNEXPECTED_SYSTEM_ERROR`. In case of these errors, contact support for assistance before continuing. 
  * `group_id` - Identifies the group where the stream is created. 
  * `group_name` - The group name where the stream is created. 
  * `properties` - List of properties associated with the stream. 
    * `property_id` - The identifier of the property. 
    * `property_name` - The descriptive label for the property. 
  * `stream_id` - A stream's unique identifier.
  * `stream_name` - The name of the stream. 
  * `stream_type_name` - Specifies the type of the data stream. `Logs - Raw` is the only stream type name currently available. 
  * `stream_version_id` - A stream version's unique identifier.
