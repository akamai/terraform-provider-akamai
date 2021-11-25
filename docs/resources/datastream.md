---
layout: "akamai"
page_title: "Akamai: DataStream"
subcategory: "DataStream"
description: |-
  DataStream
---

# akamai_datastream

Akamai constantly gathers log entries from thousands of edge servers around the world. You can use the `akamai_datastream` resource to capture these logs and deliver them to a connector of your choice at low latency. A connector, also known as a destination, represents a third-party configuration where you want to send your stream’s log files to. For each stream, you can only set one connector.

When creating a stream, you select properties to associate with the stream, data set fields to monitor in logs, and a destination to send these logs to. You can also decide whether to activate the stream on making the request. Only active streams collect and send logs to their destinations.

## Example usage

Basic usage:

```hcl
resource "akamai_datastream" "stream" {
    active             = false
    config {
        delimiter          = "SPACE"
        format             = "STRUCTURED"
        frequency {
            time_in_sec = 30
        }
        upload_file_prefix = "pre"
        upload_file_suffix = "suf"
    }
    contract_id        = "2-FGHIJ"
    dataset_fields_ids = [
        1002, 1005, 1006
    ]
    email_ids          = [
        "test@akamai.com",
        "test2@akamai.com"
    ]
    group_id           = 12345
    property_ids       = [
        100011011
    ]
    stream_name        = "Test data stream"
    stream_type        = "RAW_LOGS"
    template_name      = "EDGE_LOGS"

    s3_connector {
        access_key        = "1T2ll1H4dXWx5itGhpc7FlSbvvOvky1098nTtEMg"
        bucket            = "datastream.akamai.com"
        connector_name    = "S3Destination"
        path              = "log/edgelogs"
        region            = "ap-south-1"
        secret_access_key = "AKIA6DK7TDQLVGZ3TYP1"
    }
}
```

## Argument reference

The resource supports these arguments:

* `active` - (Required) Whether you want to start activating the stream when applying the resource. Either `true` for activating the stream upon sending the request or `false` for leaving the stream inactive after the request.
* `config` - (Required) Provides information about the log line configuration, log file format, names of log files sent, and file delivery. The argument includes these sub-arguments:
  * `delimiter` - (Optional) A delimiter that you want to use to separate data set fields in the log lines. Currently, `SPACE` is the only available delimiter. This field is required for the `STRUCTURED` log file `format`.
  * `format` - (Required) The format in which you want to receive log files, either `STRUCTURED` or `JSON`. When `delimiter` is present in the request, `STRUCTURED` is the mandatory format.
  * `frequency` - (Required) How often you want to collect logs from each uploader and send them to a destination.
      * `time_in_sec` - (Required) The time in seconds after which the system bundles log lines into a file and sends it to a destination. `30` or `60` are the possible values.
  * `upload_file_prefix` - (Optional) The prefix of the log file that you want to send to a destination. It’s a string of at most 200 characters. If unspecified, defaults to `ak`.
  * `upload_file_suffix` - (Optional) The suffix of the log file that you want to send to a destination. It’s a static string of at most 10 characters. If unspecified, defaults to `ds`.
* `contract_id` - (Required) Identifies the contract that has access to the product.
* `dataset_fields_ids` - (Required)	Identifiers of the data set fields within the template that you want to receive in logs. The order of the identifiers define how the value for these fields appears in the log lines.
* `email_ids` - (Optional) A list of email addresses you want to notify about activations and deactivations of the stream.
* `group_id` - (Required) Identifies the group that has access to the product and this stream configuration.
* `property_ids` - (Required) Identifies the properties that you want to monitor in the stream. Note that a stream can only log data for active properties.
* `stream_name` - (Required) The name of the stream.
* `stream_type` - (Required) The type of stream that you want to create. Currently, `RAW_LOGS` is the only possible stream type.
* `template_name` - (Required) The name of the data set template available for the product that you want to use in the stream. Currently, `EDGE_LOGS` is the only data set template available.
* `s3_connector` - (Optional) Specify details about the Amazon S3 connector in a stream. When validating this connector, DataStream uses the provided `access_key` and `secret_access_key` values and saves an `akamai_write_test_2147483647.txt` file in your Amazon S3 folder. You can only see this file if validation succeeds, and you have access to the Amazon S3 bucket and folder that you’re trying to send logs to. The argument includes these sub-arguments:
  * `access_key` - (Required) The access key identifier that you use to authenticate requests to your Amazon S3 account. See [Managing access keys (AWS API)](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html#Using_CreateAccessKey_API).
  * `bucket` - (Required) The name of the Amazon S3 bucket. See [Working with Amazon S3 Buckets](https://docs.aws.amazon.com/AmazonS3/latest/userguide/creating-buckets-s3.html).
  * `connector_name` - (Required) The name of the connector.
  * `path` - (Required) The path to the folder within your Amazon S3 bucket where you want to store your logs. See [Amazon S3 naming conventions](https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-keys.html#object-key-guidelines).
  * `region` - (Required) The AWS region where your Amazon S3 bucket resides. See [Regions and Zones in AWS](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html).
  * `secret_access_key` - (Required) The secret access key identifier that you use to authenticate requests to your Amazon S3 account.
* `azure_connector` - (Optional) Specify details about the Azure Storage connector configuration in a data stream. Note that currently DataStream supports only streaming data to [block objects](https://docs.microsoft.com/en-us/rest/api/storageservices/understanding-block-blobs--append-blobs--and-page-blobs). The argument includes these sub-arguments:
  * `access_key` - (Required) Either of the access keys associated with your Azure Storage account. See [View account access keys in Azure](https://docs.microsoft.com/en-us/azure/storage/common/storage-account-keys-manage?tabs=azure-portal).
  * `account_name` - (Required) Specifies the Azure Storage account name.
  * `connector_name` - (Required) The name of the connector.
  * `container_name` - (Required) Specifies the Azure Storage container name.
  * `path` - (Required) The path to the folder within the Azure Storage container where you want to store your logs. See [Azure blob naming conventions](https://docs.microsoft.com/en-us/rest/api/storageservices/naming-and-referencing-containers--blobs--and-metadata).
* `datadog_connector` - (Optional) Specify details about the Datadog connector in a stream, including:
  * `auth_token` - (Required) The API key associated with your Datadog account. See [View API keys in Datadog](https://docs.datadoghq.com/account_management/api-app-keys/#api-keys).
  * `compress logs` - (Optional) Enables GZIP compression for a log file sent to a destination. If unspecified, this defaults to `false`.
  * `connector_name` - (Required) The name of the connector.
  * `service` - (Optional) The service of the Datadog connector. A service groups together endpoints, queries, or jobs for the purposes of scaling instances. See [View Datadog reserved attribute list](https://docs.datadoghq.com/logs/log_configuration/attributes_naming_convention/#reserved-attributes).
  * `source` - (Optional) The source of the Datadog connector. See [View Datadog reserved attribute list](https://docs.datadoghq.com/logs/log_collection/?tab=http#reserved-attributes).
  * `tags` - (Optional) The tags of the Datadog connector. See [View Datadog tags](https://docs.datadoghq.com/getting_started/tagging/).
  * `url` - (Required) The Datadog endpoint where you want to store your logs. See [View Datadog logs endpoint](https://docs.datadoghq.com/logs/log_collection/?tab=http#datadog-logs-endpoints).
* `splunk_connector` - (Optional) Specify details about the Splunk connector in your stream. Note that currently DataStream supports only endpoint URLs ending with `collector/raw`. The argument includes these sub-arguments:
  * `compress_logs` - (Optional) Enables GZIP compression for a log file sent to a destination. If unspecified, this defaults to `true`.
  * `connector_name` - (Required) The name of the connector.
  * `event_collector_token` - (Required) The Event Collector token associated with your Splunk account. See [View usage of Event Collector token in Splunk](https://docs.splunk.com/Documentation/Splunk/8.0.3/Data/UsetheHTTPEventCollector).
  * `url` - (Required) The raw event Splunk URL where you want to store your logs.
* `gcs_connector` - (Optional) Specify details about the Google Cloud Storage connector you can use in a stream. When validating this connector, DataStream uses the private access key to create an `Akamai_access_verification_<timestamp>.txt` object file in your GCS bucket. You can only see this file if the validation process is successful, and you have access to the Google Cloud Storage bucket where you are trying to send logs. The argument includes these sub-arguments:
  * `bucket` - (Required) The name of the storage bucket you created in your Google Cloud account. See [Bucket naming conventions](https://cloud.google.com/storage/docs/naming-buckets).
  * `connector_name` - (Required) The name of the connector.
  * `path` - (Optional) The path to the folder within your Google Cloud bucket where you want to store logs. See [Object naming guidelines](https://cloud.google.com/storage/docs/naming-objects).
  * `private_key` - (Required) The contents of the JSON private key you generated and downloaded in your Google Cloud Storage account.
  * `project_id` - (Required) The unique ID of your Google Cloud project.
  * `service_account_name` - (Required)	The name of the service account with the storage.object.create permission or Storage Object Creator role.
* `https_connector`- (Optional) Specify details about the custom HTTPS endpoint you can use as a connector for a stream, including:
  * `authentication_type` - (Required) Either `NONE` for no authentication, or `BASIC`. For basic authentication, provide the `user_name` and `password` you set in your custom HTTPS endpoint.
  * `compress_logs` - (Optional) Whether to enable GZIP compression for a log file sent to a destination. If unspecified, this defaults to `false`.
  * `connector_name` - (Required) The name of the connector.
  * `password` - (Optional) Enter the password you set in your custom HTTPS endpoint for authentication.
  * `url` - (Required) Enter the secure URL where you want to send and store your logs.
  * `user_name` - (Optional) Enter the valid username you set in your custom HTTPS endpoint for authentication.
* `sumologic_connector` - (Optional) Specify details about the Sumo Logic connector in a stream, including:
  * `collector_code` - (Required) The unique HTTP collector code of your Sumo Logic `endpoint`.
  * `compress_logs` - (Optional)Enables GZIP compression for a log file sent to a destination. If unspecified, this defaults to `true`.
  * `connector_name` - (Required) The name of the connector.
  * `endpoint` - (Required) The Sumo Logic collection endpoint where you want to send your logs. You should follow the `https://<SumoEndpoint>/receiver/v1/http` format and pass the collector code in the `collectorCode` argument.
* `oracle_connector`- (Optional) Specify details about the Oracle Cloud Storage connector in a stream. When validating this connector, DataStream uses the provided `access_key` and `secret_access_key` values and tries to save an `Akamai_access_verification_<timestamp>.txt` file in your Oracle Cloud Storage folder. You can only see this file if the validation process is successful, and you have access to the Oracle Cloud Storage bucket and folder that you’re trying to send logs to.
  * `access_key` - (Required) The access key identifier that you use to authenticate requests to your Oracle Cloud account. See [Managing user credentials in OCS](https://docs.oracle.com/en-us/iaas/Content/Identity/Tasks/managingcredentials.htm).
  * `bucket` - (Required) The name of the Oracle Cloud Storage bucket. See [Working with Oracle Cloud Storage buckets](https://docs.oracle.com/en-us/iaas/Content/Object/Tasks/managingbuckets.htm).
  * `connector_name` - (Required) The name of the connector.
  * `namespace` - (Required) The namespace of your Oracle Cloud Storage account. See [Understanding Object Storage namespaces](https://docs.oracle.com/en-us/iaas/Content/Object/Tasks/understandingnamespaces.htm).
  * `path` - (Required) The path to the folder within your Oracle Cloud Storage bucket where you want to store your logs.
  * `region` - (Required) The Oracle Cloud Storage region where your bucket resides. See [Regions and availability domains in OCS](https://docs.oracle.com/en-us/iaas/Content/General/Concepts/regions.htm).
  * `secret_access_key` - (Required) The secret access key identifier that you use to authenticate requests to your Oracle Cloud account.

## Attributes reference

This resource returns these attributes:

* `created_by` - The user who created the stream.
* `created_date` - The date and time when the stream was created.
* `group_name` - The name of the user group that you created the stream for.
* `modified_by` - The user who modified the stream.
* `modified_date` - The date and time when the stream was modified.
* `papi_json` - The JSON-encoded rule you need to include in the property rule tree to enable the DataStream behavior. See the [DataStream workflow](../guides/get_started_datastream.md#add-a-datastream-rule-to-a-property) for more information.
* `product_id` - The ID of the product that you created stream for.
* `product_name` - The name of the product that you created this stream for.
* `stream_version_id` - Identifies the configuration version of the stream.
* `compress_logs` - Whether the GZIP compression for a log file was sent to a destination.
* `connector_id` - Identifies the connector associated with the stream.

## Import

Basic usage:

```hcl
resource "akamai_datastream" "example" {
    # (resource arguments)
  }
```

You can import your Akamai DataStream configuration using a stream version ID.

For example:

```shell
$ terraform import akamai_datastream.example 1234
```
