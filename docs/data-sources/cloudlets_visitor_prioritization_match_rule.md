---
layout: "akamai"
page_title: "Akamai: akamai_cloudlets_visitor_prioritization_match_rule"
subcategory: "Cloudlets"
description: |-
 Visitor Prioritization match rule
---

# akamai_cloudlets_visitor_prioritization_match_rule

Every policy version specifies the match rules that govern how the Cloudlet is used. Matches specify conditions that need to be met in the incoming request.

Use the `akamai_cloudlets_visitor_prioritization_match_rule` data source to build a match rule JSON object for the Visitor Prioritization Cloudlet.

## Basic usage

This example returns the JSON-encoded rules for the Visitor Prioritization Cloudlet:

```hcl
data "akamai_cloudlets_visitor_prioritization_match_rule" "example" {
    match_rules {
        name = "rule"
        start = 1644865045
        end = 1645037845
        match_url = "example.com"
        pass_through_percent = -1
        matches {
            case_sensitive = false
            match_type = "protocol"
            match_operator = "equals"
            match_value    = "http"
            check_ips = "CONNECTING_IP"
            negate = false
        }
    }
}

```

## Argument reference

This data source supports these arguments:

* `match_rules` - (Optional) A list of Cloudlet-specific match rules for a policy.
  * `name` - (Optional) The name of the rule.
  * `type` - (Optional) The type of Cloudlet the rule is for. For example, the string for Visitor Prioritization is `vpMatchRule`.
  * `start` - (Optional) The start time for this match. Specify the value in UTC in seconds since the epoch.
  * `end` - (Optional) The end time for this match. Specify the value in UTC in seconds since the epoch.
  * `matches` - (Optional) A list of conditions to apply to a Cloudlet, including:
      * `match_type` - (Optional) The type of match used, either `header`, `hostname`, `path`, `extension`, `query`, `cookie`, `deviceCharacteristics`, `clientip`, `continent`, `countrycode`, `regioncode`, `protocol`, `method`, or `proxy`.
      * `match_value` - (Optional) This depends on the `match_type`. If the `match_type` is `hostname`, then `match_value` is the fully qualified domain name, like `www.akamai.com`.
      * `match_operator` - (Optional) Compares a string expression with a pattern, either `contains`, `exists`, or `equals`.
      * `case_sensitive` - (Optional) Whether the match is case sensitive.
      * `negate` - (Optional) Whether to negate the match.
      * `check_ips` - (Optional) For `clientip`, `continent`, `countrycode`, `proxy`, and `regioncode` match types, this defines the part of the request that determines the IP address to use. Values include the connecting IP address (`CONNECTING_IP`) and the X_Forwarded_For header (`XFF_HEADERS`). To select both, enter the two values separated by a space delimiter. When both values are included, the connecting IP address is evaluated first.
      * `object_match_value` - (Optional) If `match_value` is empty, this argument is required. An object used when a rule includes more complex match criteria, like multiple value attributes. Includes these sub-arguments:
          * `name` - (Optional) If you're using a `match_type` that supports name attributes, specify the part the incoming request to match on, either `cookie`, `header`, `parameter`, or `query`.
          * `type` - (Required) The type of the array, either `object` or `simple`. Use the `simple` option when adding only an array of string-based values.
          * `name_case_sensitive` - (Optional) Whether the `name` argument should be evaluated based on case sensitivity.
          * `name_has_wildcard` - (Optional) Whether the `name` argument includes wildcards.
          * `options` - (Optional) If you set the `type` argument to `object`, use this array to list the values to match on.
              * `value` - (Optional) Specify the values in the incoming request to match on.
              * `value_has_wildcard` - (Optional) Whether the `value` argument includes wildcards.
              * `value_case_sensitive` - (Optional) Whether the `value` argument should be evaluated based on case sensitivity.
              * `value_escaped` - (Optional) Whether the `value` argument should be compared in an escaped form.
          * `value` - (Optional) If you set the `type` argument to `simple`, specify the values in the incoming request to match on.
* `match_url` - (Optional) If you're using a URL match, this specifies the URL that the Cloudlet uses to match the incoming request.
* `pass_through_percent`- (Required) Entering a value in the range of `0.0` to `99.0` specifies the percent of requests that pass through to the origin. Enter `100` to always have the request pass through to the origin. Enter `-1` to send everyone to the waiting room.
* `disabled` - (Optional) Whether to disable a rule so it is not evaluated against incoming requests. 

## Attributes reference

This data source returns these attributes:

* `type` - The type of Cloudlet the rule is for.
* `json` - A `match_rules` JSON structure generated from the API schema that defines the rules for this policy.
