---
layout: "akamai"
page_title: "Module: DataStream"
description: |-
  DataStream module for the Akamai Terraform Provider
---

# DataStream Guide

With this module, you can create data streams for your properties to provide scalable, low latency streaming of data in raw form. You can use raw data logs to find details about specific incidents, search the logs for instances using a specific IP address, or analyze the patterns of multiple attacks.

You can configure your data stream to bundle and push logs to a destination for storing, monitoring, and analytical purposes. Each data stream supports only one endpoint to send log data to, either:
* **Amazon S3.** Provides cloud object storage for your data. For more information, see [Getting started with Amazon S3](https://aws.amazon.com/s3/getting-started/).
* **Azure Storage.** Provides object storage for data objects that is highly available, secure, durable, scalable, and redundant. See [Azure Blob Storage](https://azure.microsoft.com/en-us/services/storage/blobs/).
* **Custom HTTPS endpoint.** Send log data gathered by your stream to an HTTPS endpoint of your choice.
* **Datadog.** Provides monitoring for servers, databases, and services, allows for stacking and aggregating metrics. See [Datadog](https://docs.datadoghq.com/).
* **Google Cloud Storage.** Provides a cloud-based storage with low latency, high durability, and worldwide accessibility. See [Google Cloud Storage](https://cloud.google.com/resource-manager/docs/how-to).
* **Oracle Cloud.** Provides a scalable, cloud-based storage using the S3-compatible API connectivity option to store data. See [Oracle Cloud](https://docs.oracle.com/en-us/iaas/Content/Object/home.htm).
* **Splunk.** Provides an interface for advanced metrics, monitoring, and data analysis. See [Splunk](https://docs.splunk.com/Documentation/Splunk/latest/Data/UsetheHTTPEventCollector#Send_data_to_HTTP_Event_Collector).
* **Sumo Logic.** Provides advanced analytics for your log files. See [Sumo Logic](https://help.sumologic.com/03Send-Data/Sources/02Sources-for-Hosted-Collectors/HTTP-Source).

## Prerequisites

DataStream collects performance logs against selected delivery properties and streams them to configured destinations. To create a stream, you need to have at least one existing property within the group and contract you want the stream to collect logs for. A stream can start collecting logs only if the referenced properties have the [`datastream`](https://developer.akamai.com/api/core_features/property_manager/vlatest.html#datastream) behavior enabled in their rule tree and are active on the production network.

Use the [akamai_property](../resources/property.md) and [akamai_property_activation](../resources/property_activation.md) resources to create and activate or import delivery properties.

## DataStream workflow

* [Get the group ID](#get-the-group-id). You need your group ID to create a data stream.
* [Get the data set fields](#get-the-data-set-fields). Choose the data set fields that you want to monitor in your logs in a stream configuration.
* [Create a data stream](#create-a-data-stream). Create a new stream to collect logs for associated properties and send that data to a connector.
* [Add a DataStream rule to a property](#add-a-datastream-rule-to-a-property). Copy the returned JSON snippet into the rule tree configuration.
* [Activate the property version](#activate-the-property-version). After you modify the rule tree, activate the changed property on the production network.
* [Activate the stream](#activate-the-data-stream-version) Activate the latest version of a stream to start collecting and sending logs to a destination.
* [View activation history](#view-activation-history) Check a history of activation status changes for all versions of a stream.
* [Delete a data stream](#delete-a-data-stream) If you don't need the log data anymore, you can delete a stream.

## Get the group ID

When setting up streams, you need to get the Akamai [`group_id`](../data-sources/group.md).

-> **Note** The DataStream module supports both ID formats, either with or without the `grp_` prefix. For more information about prefixes, see the [ID prefixes](https://developer.akamai.com/api/core_features/property_manager/v1.html#prefixes) section of the Property Manager API (PAPI) documentation.

## Get the data set fields

Use the [akamai_datastream_dataset_fields](../data-sources/datastream_dataset_fields.md) data source to view the data set fields available within the template. Store the `dataset_field_id` values of the fields you want to receive in logs.

## Create a data stream

To monitor and gain real-time access to delivery performance, create a new stream or import an existing one using the [akamai_datastream](../resources/datastream.md) resource. You can associate up to 100 properties with a single stream and specify a data set that you want this stream to deliver. For each property, you can create up to 3 streams to specify different data sets that you want to receive about your application, and send it to the destinations of your choice.

-> **Note** Data stream activation might be time-consuming, so set the `active` flag to `false` until you completely finish the setup.

Once you set up the `akamai_datastream` resource, run `terraform apply`. Terraform shows an overview of changes, so you can still go back and modify the configuration, or confirm to proceed. See [Command: apply](https://www.terraform.io/docs/commands/apply.html)

## Add a DataStream rule to a property

To start collecting logs for properties in a stream, you need to enable the DataStream behavior in each property that is part of any stream. You can't receive logs from properties with a disabled DataStream behavior even if they're part of active data streams.

The `terraform apply` command returns a `papi_json` output with a JSON-encoded rule for DataStream, for example:

```json
{
    "name": "Datastream Rule",
    "children": [],
    "behaviors": [
        {
            "name": "datastream",
            "options": {
                "streamType": "LOG",
                "logEnabled": true,
                "logStreamName": 7050,
                "samplingPercentage": 100
            }
        }
    ],
    "criteria": [],
    "criteriaMustSatisfy": "all"
}
```

Copy this snippet to the rule tree files in properties you created the stream for. You can also create a new `.json` file with the snippet and insert it to the property rule tree by adding `"#include:example-file.json"` under the `children` array. See [Referencing sub-files from a template](../data-sources/property_rules_template.md#referencing-sub-files-from-a-template) for more information.

If you wish to customize how your data stream is handled, see the [`datastream` behavior in the PAPI Catalog Reference](https://developer.akamai.com/api/core_features/property_manager/vlatest.html#datastream).

## Activate the property version

Use the [akamai_property_activation](../resources/property_activation.md) resource to activate the modified property version on the production network. You can only stream logs for active properties with the DataStream behavior enabled.

Run `terraform apply` again to implement the changes.

## Activate the data stream version

Once you've made all the modifications in your data stream, set the `active` flag in the `akamai_datastream` resource to `true` and run `terraform apply`. This operation takes approximately 90 minutes.

The moment a stream goes active and the DataStream behavior is enabled in your property, it starts collecting and sending logs to a destination. If you want to stop receiving these logs, you can deactivate a stream at any time by setting the flag back to `false`.

## View activation history

Use the [akamai_datastream_activation_history](../data-sources/datastream_activation_history.md) data source to get detailed information about activation status changes for a version of a stream.

## Delete a stream

To delete a stream, remove the `akamai_datastream` resource and all the dependencies from your Terraform configuration. If you want delete an `active` stream, the provider automatically deactivates it first. If you want to delete a stream with a pending status, either `activating` or `deactivating`, the provider waits until the status becomes stable and proceeds with the operation.

Deleting a stream means that you can’t activate this stream again, and that you stop receiving logs for the properties that this stream monitors.
