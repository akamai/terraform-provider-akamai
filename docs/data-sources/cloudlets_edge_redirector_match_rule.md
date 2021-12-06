---
layout: "akamai"
page_title: "Akamai: akamai_cloudlets_edge_redirector_match_rule"
subcategory: "Cloudlets"
description: |-
 Edge Redirector match rule
---

# akamai_cloudlets_edge_redirector_match_rule

Every policy version specifies the match rules that govern how the Cloudlet is used. Matches specify conditions that need to be met in the incoming request.

Use the `akamai_cloudlets_edge_redirector_match_rule` data source to build a match rule JSON object for the Edge Redirector Cloudlet.

## Basic usage

This example returns the JSON-encoded rules for the Edge Redirector Cloudlet:

```hcl
data "akamai_cloudlets_edge_redirector_match_rule" "example" {
    match_rules {
        name = "rule"
        start = 10
        end = 10000
        match_url = "example.com"
        redirect_url = "https://www.example.com"
        status_code = 301
        use_incoming_query_string = false
        use_relative_url = "none"
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
    }
}
```

## Argument reference

This data source supports these arguments:

* `match_rules` - (Optional) A list of Cloudlet-specific match rules for a policy.
  * `name` - (Optional) The name of the rule.
  * `type` - (Optional) The type of Cloudlet the rule is for. For example, the string for Edge Redirector is `erMatchRule`.
  * `start` - (Optional) The start time for this match. Specify the value in UTC in seconds since the epoch.
  * `end` - (Optional) The end time for this match. Specify the value in UTC in seconds since the epoch.
  * `matches` - (Optional) A list of conditions to apply to a Cloudlet, including:
      * `match_type` - (Optional) The type of match used, either `header`, `hostname`, `path`, `extension`, `query`, `regex`, `cookie`, `deviceCharacteristics`, `clientip`, `continent`, `countrycode`, `regioncode`, `protocol`, `method`, or `proxy`.
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
* `use_relative_url` - (Optional) If set to `relative_url`, takes the path you specify in the `redirect_url` argument and sets it in the response’s Location header. The client or browser receiving the request decides which protocol and hostname to use. If set to `copy_scheme_hostname`, creates an absolute path by taking the protocol and hostname from the incoming request and combining them with path information you specify in the `redirect_url` argument. This absolute path is set in the response’s Location header. If you do not specify use_relative_url or set to `none`, then specify the `redirect_url` argument as a fully-qualified URL.
* `status_code` - (Required) The HTTP response status code, which is either `301` (permanent redirect) or `302` (temporary redirect).
* `redirect_url` - (Required) The URL Edge Redirector redirects the request to. If you're using `use_relative_url`, you can enter a path for the value.
* `match_url` - (Optional) If you're using a URL match, this specifies the URL that the Cloudlet uses to match the incoming request.
* `use_incoming_query_string` - (Optional) Whether the Cloudlet should include the query string from the request in the rewritten or forwarded URL.

## Attributes reference

This data source returns these attributes:

* `type` - The type of Cloudlet the rule is for.
* `json` - A `match_rules` JSON structure generated from the API schema that defines the rules for this policy.
