---
layout: "akamai"
page_title: "Akamai: cloudlets_application_load_balancer_match_rule"
subcategory: "Cloudlets"
description: |-
 Application Load Balancer match rule
---

# akamai_cloudlets_application_load_balancer_match_rule

Every policy version specifies the match rules that govern how the Cloudlet is used. Matches specify conditions that need to be met in the incoming request.

Use the `akamai_cloudlets_application_load_balancer_match_rule` data source to build a match rule JSON object for the Application Load Balancer Cloudlet.

## Basic usage

This example returns the JSON-encoded rules for the Application Load Balancer Cloudlet:

```hcl
data "akamai_cloudlets_application_load_balancer_match_rule" "example" {
    match_rules {
        name = "rule"
        start = 1644865045
        end = 1645037845
        match_url = "example.com"
        matches {
            case_sensitive = false
            match_operator = "equals"
            match_type = "method"
            negate = false
            object_match_value {
                type = "simple"
                value = ["GET"]
            }
        }
        forward_settings {
            origin_id = "alb_test_1"
        }
    }
}
```

## Argument reference

This data source supports these arguments:

* `match_rules` - (Optional) A list of Cloudlet-specific match rules for a policy.
  * `name` - (Optional) The name of the rule.
  * `type` - (Optional) The type of Cloudlet the rule is for. For example, the string for Application Load Balancer is `albMatchRule`.
  * `start` - (Optional) The start time for this match. Specify the value in UTC in seconds since the epoch.
  * `end` - (Optional) The end time for this match. Specify the value in UTC in seconds since the epoch.
  * `id` - (Optional) An identifier for Akamai internal use only.
  * `matches` - (Optional) A list of conditions to apply to a Cloudlet, including:
      * `match_type` - (Optional) The type of match used, either `clientip`, `continent`, `cookie`, `countrycode`, `deviceCharacteristics`, `extension`, `header`, `hostname`, `method`, `path`, `protocol`, `proxy`, `query`, `regioncode`, or `range`.
      * `match_value` - (Optional) This depends on the `match_type`. If the `match_type` is `hostname`, then `match_value` is the fully qualified domain name, like `www.akamai.com`.
      * `match_operator` - (Optional) Compares a string expression with a pattern, either `contains`, `exists`, or `equals`.
      * `case_sensitive` - (Optional) Whether the match is case sensitive.
      * `negate` - (Optional) Whether to negate the match.
      * `check_ips` - (Optional) For `clientip`, `continent`, `countrycode`, `proxy`, and `regioncode` match types, this defines the part of the request that determines the IP address to use. Values include the connecting IP address (`CONNECTING_IP`) and the X_Forwarded_For header (`XFF_HEADERS`). To select both, enter the two values separated by a space delimiter. When both values are included, the connecting IP address is evaluated first.
      * `object_match_value` - (Optional) If `match_value` is empty, this argument is required. An object used when a rule either includes more complex match criteria, like multiple value attributes, or a range match. Includes these sub-arguments:
          * `name` - (Optional) If you're using a `match_type` that supports name attributes, specify the part the incoming request to match on, either `cookie`, `header`, `parameter`, or `query`.
          * `type` - (Required) The type of the array, either `object`, `range`, or `simple`. Use the `simple` option when adding only an array of string-based values.
          * `name_case_sensitive` - (Optional) Whether the `name` argument should be evaluated based on case sensitivity.
          * `name_has_wildcard` - (Optional) Whether the `name` argument includes wildcards.
          * `options` - (Optional) If you set the `type` argument to `object`, use this array to list the values to match on.
              * `value` - (Optional) Specify the values in the incoming request to match on.
              * `value_has_wildcard` - (Optional) Whether the `value` argument includes wildcards.
              * `value_case_sensitive` - (Optional) Whether the `value` argument should be evaluated based on case sensitivity.
              * `value_escaped` - (Optional) Whether the `value` argument should be compared in an escaped form.
         * `value` - (Optional) If you set the `type` argument to `simple` or `range`, specify the values in the incoming request to match on. With `range`, you can only specify an array of integers, for example `[1, 2]`.
  * `match_url` - (Optional) The URL that the Cloudlet uses to match the incoming request.
  * `matches_always` - (Optional) Whether the match supports default rules that apply to all requests.
  * `forward_settings` - (Required) Defines data used to construct a new request URL if all conditions are met. If all of the conditions you set are true, the Edge Server returns an HTTP response from the rewritten URL.
    * `origin_id` - (Required) The ID of the Conditional Origin the requests are forwarded to.

## Attributes reference

This data source returns these attributes:

* `type` - The type of Cloudlet the rule is for.
* `json` - A `match_rules` JSON structure generated from the API schema that defines the rules for this policy.
