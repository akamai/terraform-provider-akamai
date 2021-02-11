---
layout: "akamai"
page_title: "Akamai: gtm property"
subcategory: "Global Traffic Management"
description: |-
  GTM Property
---

# akamai_gtm_property

Use the `akamai_gtm_property` resource provides the resource for creating, configuring and importing a GTM property, a set of IP addresses or CNAMEs that GTM provides in response to DNS queries based on a set of rules. 

~> **Note** Import requires an ID with this format: `existing_domain_name`:`existing_property_name`.

## Example usage

Basic usage:

```hcl
resource "akamai_gtm_property" "demo_property" {
    domain = "demo_domain.akadns.net"
    name = "demo_property"
    type =  "weighted-round-robin"
    score_aggregation_type = "median"
    handout_limit = 5
    handout_mode = "normal"
    traffic_target {
        datacenter_id = 3131
    }
}
```

## Argument reference

This resource supports these arguments:

* `domain` - (Required) DNS name for the GTM Domain set that includes this Property.
* `name` - (Required) DNS name for a collection of IP address or CNAME responses. The value, together with the GTM domainName, forms the Property’s hostname. 
* `type` - (Required) Specifies the load balancing behavior for the property. Either failover, geographic, cidrmapping, weighted-round-robin, weighted-hashed, weighted-round-robin-load-feedback, qtr, or performance. 
* `score_aggregation_type` - (Required) Specifies how GTM aggregates liveness test scores across different tests, when multiple tests are configured.
* `handout_limit` - (Required) Indicates the limit for the number of live IPs handed out to a DNS request.
* `handout_mode` - (Required) Specifies how IPs are returned when more than one IP is alive and available.
* `traffic_target` - (Required) Contains information about where to direct data center traffic. You can have multiple `traffic_target` arguments. If used, requires these arguments:
  * `datacenter_id` - (Required) A unique identifier for an existing data center in the domain.
  * `enabled` - (Required) A boolean indicating whether the traffic target is used. You can also omit the traffic target, which has the same result as the false value.
  * `weight` - (Required) Specifies the traffic weight for the target.
  * `servers` - (Required) (List) Identifies the IP address or the hostnames of the servers.
  * `name` - (Required) An alternative label for the traffic target.
  * `handout_cname` - (Required) Specifies an optional data center for the property. Used when there are no servers configured for the property.
* `liveness_test` - (Optional) Contains information about the liveness tests, which are run periodically to determine whether your servers respond to requests. You can have multiple `liveness_test` arguments. If used, requires these arguments:
  * `name` - (Optional) A descriptive name for the liveness test.
  * `test_interval` - (Optional) Indicates the interval at which the liveness test is run, in seconds. Requires a minimum of 10 seconds.
  * `test_object_protocol` - (Optional) Specifies the test protocol. Possible values include `DNS`, `HTTP`, `HTTPS`, `FTP`, `POP`, `POPS`, `SMTP`, `SMTPS`, `TCP`, or `TCPS`.
  * `test_timeout` - (Optional) Specifies the duration of the liveness test before it fails. The range is from 0.001 to 60 seconds.
  * `answers_required` - (Optional) If `test_object_protocol` is DNS, enter a boolean value if an answer is needed for the DNS query to be successful.
  * `disabled` - (Optional) A boolean indicating whether the liveness test is disabled. When disabled, GTM stops running the test, effectively treating it as if it no longer exists.
  * `disable_nonstandard_port_warning` - (Optional) A boolean that if set to `true`, disables warnings when non-standard ports are used.
  * `error_penalty` - (Optional) Specifies the score that’s reported if the liveness test encounters an error other than timeout, such as connection refused, and 404.
  * `http_header` - (Optional) Contains HTTP headers to send if the `test_object_protocol` is `http` or `https`. You can have multiple `http_header` entries. Requires these arguments: 
    * `name` - Name of HTTP header.
    * `value` - Value of HTTP header.
  * `http_error3xx` - (Optional) A boolean that if set to `true`, treats a 3xx HTTP response as a failure if the `test_object_protocol` is `http`, `https`, or `ftp`.
  * `http_error4xx` - (Optional) A boolean that if set to `true`, treats a 4xx HTTP response as a failure if the `test_object_protocol` is `http`, `https`, or `ftp`.
  * `http_error5xx` - (Optional) A boolean that if set to `true`, treats a 5xx HTTP response as a failure if the `test_object_protocol` is `http`, `https`, or `ftp`.
  * `peer_certificate_verification` - (Optional) A boolean that if set to `true`, validates the origin certificate. Applies only to tests with `test_object_protocol` of https.
  * `recursion_requested` - (Optional) A boolean indicating whether the `test_object_protocol` is DNS. The DNS query is recursive.
  * `request_string` - (Optional) Specifies a request string.
  * `resource_type` - (Optional) Specifies the query type, if `test_object_protocol` is DNS.
  * `response_string` - (Optional) Specifies a response string.
  * `ssl_client_certificate` - (Optional) Indicates a Base64-encoded certificate. SSL client certificates are available for livenessTests that use secure protocols.
  * `ssl_client_private_key` - (Optional) Indicates a Base64-encoded private key. The private key used to generate or request a certificate for livenessTests can’t have a passphrase nor be used for any other purpose.
  * `test_object` - (Optional) Specifies the static text that acts as a stand-in for the data that you’re sending on the network.
  * `test_object_password` - (Optional) Specifies the test object’s password. It is required if testObjectProtocol is ftp.
  * `test_object_port` - (Optional) Specifies the port number for the testObject.
  * `test_object_username` - (Optional) A descriptive name for the testObject.
  * `timeout_penalty`- (Optional) Specifies the score to be reported if the liveness test times out.
* `wait_on_complete` - (Optional) A boolean indicating whether to wait for transaction to complete. Set to `true` by default.
* `failover_delay` - (Optional) Specifies the failover delay in seconds.
* `failback_delay` - (Optional) Specifies the failback delay in seconds.
* `ipv6` - (Optional) A boolean that indicates the type of IP address handed out by a GTM property.
* `stickiness_bonus_percentage` - (Optional) Specifies a percentage used to configure data center affinity.
* `stickiness_bonus_constant` - (Optional) Specifies a constant used to configure data center affinity.
* `health_threshold` - (Optional) Configures a cutoff value that is computed from the median scores.
* `use_computed_targets` - (Optional) For load-feedback domains only, a boolean that indicates whether you want GTM to automatically compute target load.
* `backup_ip` - Specifies a backup IP. When GTM declares that all of the targets are down, the backupIP is handed out.
* `balance_by_download_score` - (Optional) A boolean that indicates whether download score based load balancing is enabled.
* `static_ttl` - (Optional) Specifies the TTL in seconds for static resource records that do not change based on the requesting name server IP.
* `unreachable_threshold` - (Optional) For performance domains, this specifies a penalty value that’s added to liveness test scores when data centers have an aggregated loss fraction higher than this value.
* `health_multiplier` - (Optional) Configures a cutoff value that is computed from the median scores.
* `dynamic_ttl` - (Optional) Indicates the TTL in seconds for records that might change dynamically based on liveness and load balancing such as A and AAAA records, and CNAMEs.
* `max_unreachable_penalty` - (Optional) For performance domains, this specifies a penalty value that’s added to liveness test scores when data centers show an aggregated loss fraction higher than the penalty value.
* `map_name` - (Optional) A descriptive label for a GeographicMap or a CidrMap that’s required if the property is either geographic or cidrmapping, in which case mapName needs to reference either an existing GeographicMap or CidrMap in the same domain.
* `load_imbalance_percentage` - (Optional) Indicates the percent of load imbalance factor (LIF) for the property.
* `health_max` - (Optional) Defines the absolute limit beyond which IPs are declared unhealthy.
* `cname` - (Optional) Indicates the fully qualified name aliased to a particular property.
* `comments` - (Optional) A descriptive note about changes to the domain. The maximum is 4000 characters.
* `ghost_demand_reporting` - (Optional) Use load estimates from Akamai Ghost utilization messages.
* `min_live_fraction` - (Optional) Specifies what fraction of the servers need to respond to requests so GTM considers the data center up and able to receive traffic.
* `static_rr_set` - (Optional) Contains static record sets. You can have multiple `static_rr_set` entries. Requires these arguments: 
  * `type` - (Optional) The record type.
  * `ttl` - (Optional) The number of seconds that this record should live in a resolver’s cache before being refetched.
  * `rdata` - (Optional) (List) An array of data strings, representing multiple records within a set.

### Computed arguments

This resource returns these computed arguments in the `terraform.tfstate` file:

* `weighted_hash_bits_for_ipv4`
* `weighted_hash_bits_for_ipv6`

### Schema reference

You can download the GTM Property backing schema from the [Global Traffic Management API](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#property) page.
