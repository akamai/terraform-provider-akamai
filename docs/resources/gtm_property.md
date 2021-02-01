---
layout: "akamai"
page_title: "Akamai: gtm property"
subcategory: "Global Traffic Management"
description: |-
  GTM Property
---

# akamai_gtm_property

`akamai_gtm_property` provides the resource for creating, configuring and importing a GTM property, a set of IP addresses or CNAMEs that GTM provides in response to DNS queries based on a set of rules. Note: Import requires an ID of the format: `existing_domain_name`:`existing_property_name`.

## Example Usage

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

## Argument Reference

The following arguments are supported:

### Required

* `domain` - DNS name for the GTM Domain set that includes this Property.
* `name` - DNS name for a collection of IP address or CNAME responses. The value, together with the GTM domainName, forms the Property’s hostname. 
* `type` - Specifies the load balancing behavior for the property. Either failover, geographic, cidrmapping, weighted-round-robin, weighted-hashed, weighted-round-robin-load-feedback, qtr, or performance. 
* `score_aggregation_type` - Specifies how GTM aggregates liveness test scores across different tests, when multiple tests are configured.
* `handout_limit` - Indicates the limit for the number of live IPs handed out to a DNS request.
* `handout_mode` - Specifies how IPs are returned when more than one IP is alive and available.
* `traffic_target` - (multiple allowed) Contains information about where to direct data center traffic.
  * `datacenter_id` - A unique identifier for an existing data center in the domain.
  * `enabled` - (Boolean) Indicates whether the traffic target is used. You can also omit the traffic target, which has the same result as the false value.
  * `weight` - Specifies the traffic weight for the target.
  * `servers` - (List) Identifies the IP address or the hostnames of the servers.
  * `name` - An alternative label for the traffic target.
  * `handout_cname` - Specifies an optional data center for the property. The handoutCName is used when there are no servers configured for the property.

### Optional

* `liveness_test` - (multiple allowed) Contains information about the liveness tests, which are run periodically to determine whether your servers respond to requests.
  * `name` - A descriptive name for the liveness test.
  * `test_interval` - Indicates the interval at which the liveness test is run, in seconds. Requires a minimum of 10 seconds.
  * `test_object_protocol` - Specifies the test protocol. EitherDNS, HTTP, HTTPS, FTP, POP, POPS, SMTP, SMTPS, TCP, or TCPS.
  * `test_timeout` - Specifies the duration of the liveness test before it fails. The range is from 0.001 to 60 seconds.
  * `answers_required` - (Boolean) If testObjectProtocol is DNS, requires an answer to the DNS query to be considered a success.
  * `disabled` - (Boolean) Disables the liveness test. When disabled, GTM will cease to run the test, effectively treating it as if it no longer exists.
  * `disable_nonstandard_port_warning` - (Boolean) Disables warnings when non-standard ports are used.
  * `error_penalty` - Specifies the score that’s reported if the liveness test encounters an error other than timeout, such as connection refused, and 404.
  * `http_header` - (multiple allowed) Contains HTTP headers to send if the testObjectProtocol is http or https.
     `name` - Name of HTTP header.
     `value` - Value of HTTP header.
  * `http_error3xx` - (Boolean) Treats a 3xx HTTP response as a failure if the testObjectProtocol is http, https, or ftp.
  * `http_error4xx` - (Boolean) Treats a 4xx HTTP response as a failure if the testObjectProtocol is http, https, or ftp.
  * `http_error5xx` - (Boolean) Treats a 5xx HTTP response as a failure if the testObjectProtocol is http, https, or ftp.
  * `peer_certificate_verification` - (Boolean) Validates the origin certificate. Applies only to tests with testObjectProtocol of https.
  * `recursion_requested` - (Boolean) Indicates that testObjectProtocol is DNS. The DNS query is recursive.
  * `request_string` - Specifies a request string.
  * `resource_type` - Specifies the query type, if testObjectProtocol is DNS.
  * `response_string` - Specifies a response string.
  * `ssl_client_certificate` - Indicates a Base64-encoded certificate. SSL client certificates are available for livenessTests that use secure protocols.
  * `ssl_client_private_key` - Indicates a Base64-encoded private key. The private key used to generate or request a certificate for livenessTests can’t have a passphrase nor be used for any other purpose.
  * `test_object` - Specifies the static text that acts as a stand-in for the data that you’re sending on the network.
  * `test_object_password` - Specifies the test object’s password. It is required if testObjectProtocol is ftp.
  * `test_object_port` - Specifies the port number for the testObject.
  * `test_object_username` - A descriptive name for the testObject.
  * `timeout_penalty`- Specifies the score to be reported if the liveness test times out.
* `wait_on_complete` - (Boolean, Default: true) Wait for transaction to complete.
* `failover_delay` - Specifies the failover delay in seconds.
* `failback_delay` - Specifies the failback delay in seconds.
* `ipv6` - (Boolean) Indicates the type of IP address handed out by a Property.
* `stickiness_bonus_percentage` - Specifies a percentage used to configure data center affinity.
* `stickiness_bonus_constant` - Specifies a constant used to configure data center affinity.
* `health_threshold` - Configures a cutoff value that is computed from the median scores.
* `use_computed_targets` - (Boolean) For load-feedback domains only, this specifies that you want GTM to automatically compute target load.
* `backup_ip` - Specifies a backup IP. When GTM declares that all of the targets are down, the backupIP is handed out.
* `balance_by_download_score` - (Boolean) Indicates whether download score based load balancing is enabled.
* `static_ttl` - Specifies the TTL in seconds for static resource records that do not change based on the requesting name server IP.
* `unreachable_threshold` - For performance domains, this specifies a penalty value that’s added to liveness test scores when data centers have an aggregated loss fraction higher than this value.
* `health_multiplier` - Configures a cutoff value that is computed from the median scores.
* `dynamic_ttl` - Indicates the TTL in seconds for records that might change dynamically based on liveness and load balancing such as A and AAAA records, and CNAMEs.
* `max_unreachable_penalty` - For performance domains, this specifies a penalty value that’s added to liveness test scores when data centers show an aggregated loss fraction higher than the penalty value.
* `map_name` - A descriptive label for a GeographicMap or a CidrMap that’s required if the property is either geographic or cidrmapping, in which case mapName needs to reference either an existing GeographicMap or CidrMap in the same domain.
* `load_imbalance_percentage` - Indicates the percent of load imbalance factor (LIF) for the property.
* `health_max` - Defines the absolute limit beyond which IPs are declared unhealthy.
* `cname` - Indicates the fully qualified name aliased to a particular property.
* `comments` - A descriptive note about changes to the domain. The maximum is 4000 characters.
* `ghost_demand_reporting` - Use load estimates from Akamai Ghost utilization messages.
* `min_live_fraction` - Specifies what fraction of the servers need to respond to requests so GTM considers the data center up and able to receive traffic.
* `static_rr_set` - (multiple allowed) Contains static recordsets.
  * `type` - The record type.
  * `ttl` - The number of seconds that this record should live in a resolver’s cache before being refetched.
  * `rdata` - (List) An array of data strings, representing multiple records within a set.

Computed

The following arguments will be found in `terraform.tfstate` and can be referenced throughout the configuration. The values cannot be changed.

* `weighted_hash_bits_for_ipv4`.
* `weighted_hash_bits_for_ipv6`.

### Schema Reference

The GTM Property backing schema and more complete element descriptions can be found at [Akamai Developer Website](https://developer.akamai.com/api/web_performance/global_traffic_management/v1.html#property).
