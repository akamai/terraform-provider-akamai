---
layout: akamai
subcategory: Cloudlets
---

# akamai_cloudlets_application_load_balancer

Use the `akamai_cloudlets_application_load_balancer` resource to create the Application Load Balancer Cloudlet configuration. The Application Load Balancer Cloudlet provides intelligent, scalable traffic management across physical, virtual, and cloud-hosted data centers without requiring the origin to send load feedback. This Cloudlet can automatically detect load conditions and route traffic to the optimal data source while maintaining custom routing policies and consistent visitor session behavior for your visitors.

## Example usage

Basic usage:

```hcl
resource "akamai_cloudlets_application_load_balancer" "example" {
  origin_id      = "alb_test_1"
  description    = "application_load_balancer description"
  balancing_type = "WEIGHTED"
  data_centers {
    cloud_server_host_header_override = false
    cloud_service = true
    country  = "US"
    continent = "NA"
    latitude = 102.78108
    longitude = -116.07064
    percent = 100
    liveness_hosts = ["example"]
    hostname= "example-hostname"
    state_or_province = "MA"
    city = "Boston"
    origin_id = "alb_test_1"
  }
  liveness_settings {
    port = 1234
    protocol = "HTTP"
    path = "/status"
    host_header = "header"
    additional_headers = {
      additional_headers = "123"
    }
    interval = 10
    request_string = "test_request_string"
    response_string = "test_response_string"
    timeout = 60
  }
}
```

## Argument reference

The following arguments are supported:

* `origin_id` - (Required) A unique identifier for the Conditional Origin that supports the load balancing configuration. The Conditional Origin type must be set to `APPLICATION_LOAD_BALANCER` in the `origin` behavior. See [property rules](../data-sources/property_rules.md) for more information.
* `description` - (Optional) The description of the load balancing configuration.
* `balancing_type` - (Optional) The type of load balancing being performed, either `WEIGHTED` or `PERFORMANCE`.
* `data_centers` - (Required) Specifies the Conditional Origins being used as data centers for an Application Load Balancer implementation. Only Conditional Origins with an origin type of `CUSTOMER` or `NETSTORAGE` can be used as data centers in an Application Load Balancer configuration.
  * `latitude` - (Required) The latitude value for the data center. This member supports six decimal places of precision.
  * `longitude` - (Required) The longitude value for the data center. This member supports six decimal places of precision.
  * `continent` - (Required) The code of the continent on which the data center is located. See [Continent Codes](https://control.akamai.com/dl/edgescape/continentCodes.csv) for a list of valid codes.
  * `country` - (Required) The country in which the data center is located. See [Country Codes](https://control.akamai.com/dl/edgescape/cc2continent.csv) for a list of valid codes.
  * `origin_id` - (Required) The identifier of an origin that represents the data center. The Conditional Origin, which is defined in Property Manager, must have an origin type of either `CUSTOMER` or `NET_STORAGE` set in the `origin` behavior. See [property rules](../data-sources/property_rules.md) for more information.
  * `percent`  - (Required) The percent of traffic that is sent to the data center. The total for all data centers must equal 100%.
  * `cloud_service` - (Optional) Whether this datacenter is a cloud service.
  * `liveness_hosts` - (Optional) A list of the origin servers used to poll the data centers in an Application Load Balancer configuration. These servers support basic HTTP polling.
  * `hostname` - (Optional) The name of the host that can be used as a Conditional Origin. This should match the `hostname` value defined for this datacenter in Property Manager.
  * `state_or_province` - (Optional) The state, province, or region where the data center is located.
  * `city` - (Optional) The city in which the data center is located.
  * `cloud_server_host_header_override` - (Optional) Whether to override the cloud server host header.
* `liveness_settings` - (Optional) Specifies the health of each load balanced data center defined in the data center list.
  * `port` - (Required) The port for the test object. The default port is 80, which is standard for HTTP. Enter 443 if you are using HTTPS.
  * `protocol` - (Required) The protocol or scheme for the database, either `HTTP` or `HTTPS`.
  * `path` - (Required) The path to the test object used for liveness testing. The function of the test object is to help determine whether the data center is functioning.
  * `host_header` - (Optional) The Host header for the liveness HTTP request.
  * `additional_headers` - (Optional) Maps additional case-insensitive HTTP header names included to the liveness testing requests.
  * `interval` - (Optional) The frequency of liveness tests. Defaults to 60 seconds, minimum is 10 seconds.
  * `peer_certificate_verification` - (Optional) Whether to validate the origin certificate for an HTTPS request.
  * `request_string` - (Optional) The request used for TCP and TCPS tests.
  * `response_string` - (Optional) The response used for TCP and TCPS tests.
  * `status_3xx_failure` - (Optional) If set to `true`, marks the liveness test as failed when the request returns a 3xx (redirection) status code.
  * `status_4xx_failure` - (Optional) If set to `true`, marks the liveness test as failed when the request returns a 4xx (client error) status code.
  * `status_5xx_failure` - (Optional) If set to `true`, marks the liveness test as failed when the request returns a 5xx (server error) status code.
  * `timeout` - (Optional) The number of seconds the system waits before failing the liveness test.

## Attribute reference

The following attributes are returned:

* `version` - The version number of the load balancing configuration.
* `warnings` - A list of warnings that occurred during the activation of the load balancing configuration.

## Import

Basic usage:

```hcl
resource "akamai_cloudlets_application_load_balancer" "example" {
    # (resource arguments)
  }
```

You can import your Akamai Application Load Balancer configuration using an origin ID.

For example:

```shell
$ terraform import akamai_datastream.example alb_test_1
```
