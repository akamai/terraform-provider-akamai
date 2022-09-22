---
layout: akamai
subcategory: Cloudlets
---

# akamai_cloudlets_application_load_balancer

Use the `akamai_cloudlets_application_load_balancer` data source to list details about the Application Load Balancer configuration with a specified policy version, or latest if not specified.

## Basic usage

This example returns the load balancing configuration details based on the origin ID and optionally, a version:

```hcl
data "akamai_cloudlets_application_load_balancer" "example" {
  origin_id = "alb_test_1"
  version = 1
}
```

## Argument reference

This data source supports these arguments:

* `origin_id` - (Required) A unique identifier for the Conditional Origin that supports the load balancing configuration. The Conditional Origin type must be set to `APPLICATION_LOAD_BALANCER` in the `origin` behavior. See [property rules](../data-sources/property_rules.md) for more information.
* `version` - (Optional) The version number of the load balancing configuration.

## Attributes reference

This data source returns these attributes:

* `description` - The description of the load balancing configuration.
* `type` - The type of Conditional Origin. `APPLICATION_LOAD_BALANCER` is the only supported value.
* `balancing_type` - The type of load balancing being performed, either `WEIGHTED` or `PERFORMANCE`.
* `created_by` - The name of the user who created this load balancing configuration.
* `created_date` - The date, in ISO 8601 format, when this load balancing configuration was created.
* `deleted` - Whether the Conditional Origin version has been deleted. If `false`, you can use this version again.
* `immutable` - Whether you can edit the load balancing version. The default setting for this member is false. It automatically becomes true when the load balancing version is activated for the first time.
* `last_modified_by` - The user who last modified the load balancing configuration.
* `last_modified_date` - The date, in ISO 8601 format, when the initial load balancing configuration was last modified.
* `warnings` - A list of warnings that occured during the activation of the load balancing configuration.
* `data_centers` - Specifies the Conditional Origins being used as data centers for an Application Load Balancer implementation. Only Conditional Origins with an origin type of `CUSTOMER` or `NETSTORAGE` can be used as data centers in an Application Load Balancer configuration.
  * `city` - The city in which the data center is located.
  * `cloud_server_host_header_override` - Whether the cloud server host header is overridden.
  * `cloud_service` - Whether this datacenter is a cloud service.
  * `continent` - The code of the continent on which the data center is located. See [Continent Codes](https://control.akamai.com/dl/edgescape/continentCodes.csv) for a list of valid codes.
  * `country` - The country in which the data center is located. See [Country Codes](https://control.akamai.com/dl/edgescape/cc2continent.csv) for a list of valid codes.
  * `hostname` - The name of the host that can be used as a Conditional Origin. This should match the `hostname` value defined for this datacenter in Property Manager.
  * `latitude` - The latitude value for the data center. This member supports six decimal places of precision.
  * `liveness_hosts` - A list of the origin servers used to poll the data centers in an Application Load Balancer configuration. These servers support basic HTTP polling.
  * `longitude` - The longitude value for the data center. This member supports six decimal places of precision.
  * `origin_id` - The ID of an origin that represents the data center. The Conditional Origin, which is defined in Property Manager, must have an origin type of either `CUSTOMER` or `NET_STORAGE` set in the `origin` behavior. See [property rules](../data-sources/property_rules.md) for more information.
  * `percent` - The percent of traffic that is sent to the data center. The total for all data centers must equal 100%.
  * `state_or_province` - The state, province, or region where the data center is located.
* `liveness_settings` - Specifies the health of each load balanced data center defined in the data center list.
  * `host_header` - The Host header for the liveness HTTP request.
  * `additional_headers` - Maps additional case-insensitive HTTP header names included to the liveness testing requests.
  * `interval` - The frequency of liveness tests. Defaults to 60 seconds, minimum is 10 seconds.
  * `path` - The path to the test object used for liveness testing. The function of the test object is to help determine whether the data center is functioning.
  * `peer_certificate_verification` - Whether to validate the origin certificate for an HTTPS request.
  * `port` - The port for the test object. The default port is 80, which is standard for HTTP. Enter 443 if you are using HTTPS.
  * `protocol` - The protocol or scheme for the database, either `HTTP` or `HTTPS`.
  * `request_string` - The request used for TCP and TCPS tests.
  * `response_string` - The response used for TCP and TCPS tests.
  * `status_3xx_failure` - If `true`, marks the liveness test as failed when the request returns a 3xx (redirection) status code.
  * `status_4xx_failure` - If `true`, marks the liveness test as failed when the request returns a 4xx (client error) status code.
  * `status_5xx_failure` - If `true`, marks the liveness test as failed when the request returns a 5xx (server error) status code.
  * `timeout` - The number of seconds the system waits before failing the liveness test.
