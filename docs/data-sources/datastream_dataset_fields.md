---
layout: "akamai"
page_title: "Akamai: DataStream Data Set Fields"
subcategory: "DataStream"
description: |-
 Data Set Fields
---

# akamai_datastream_dataset_fields

Use the `akamai_datastream_dataset_fields` data source to list groups of data set fields available in the template.

## Example usage

This example returns data set fields for a default template:

```hcl
data "akamai_datastream_dataset_fields" "fields" {}
```

## Argument reference

The data source supports this argument:

* `template_name` - (Optional) The name of the data set template you use in your stream configuration. Currently, `EDGE_LOGS` is the only available data set template and the default value for this argument.

## Attributes reference

This data source returns these attributes:

* `fields` - A group of data set fields available in a template, including:
  * `dataset_group_name` - The name of the data set group.
  * `dataset_group_description` - Additional information about the data set group.
  * `dataset_fields` - A list of data set fields available within the data set group, including:
    * `dataset_field_description` - Additional information about the data set field.
    * `dataset_field_id` - Unique identifier for the field.
    * `dataset_field_json_key` - The JSON key for the field in a log line.
    * `dataset_field_name` - The name of the data set field.
