---
layout: akamai
subcategory: Cloudlets
---

# akamai_cloudlets_policy

Use the `akamai_cloudlets_policy` resource to create and version a policy. For each Cloudlet instance on your contract, there can be any number of policies. A single policy is associated with a single property configuration. Within a policy version you define the rules that determine when the Cloudlet executes. You may want to create a new version of a policy to support a different business requirement, or to test new functionality.


## Example usage

Basic usage:

```hcl
resource "akamai_cloudlets_policy" "example" {
  name          = "policy1"
  cloudlet_code = "ER"
  description   = "policy description"
  group_id      = "grp_123"
  match_rules   = <<-EOF
  [
  {
    "name": "rule1",
    "type": "erMatchRule",
    "useRelativeUrl": "none",
    "statusCode": 301,
    "redirectURL": "https://www.example.com",
    "matchURL": "example.com",
    "useIncomingQueryString": false,
    "useIncomingSchemeAndHost": false
  },
  {
    "name": "rule2",
    "type": "erMatchRule",
    "matches": [
      {
        "matchType": "path",
        "matchValue": "/example/website.html",
        "matchOperator": "equals",
        "caseSensitive": false,
        "negate": false
      }
    ],
    "useRelativeUrl": "copy_scheme_hostname",
    "statusCode": 301,
    "redirectURL": "/website.html",
    "useIncomingQueryString": false,
    "useIncomingSchemeAndHost": true
  }
]
EOF
}
```

## Argument reference

The following arguments are supported:

* `name` - (Required) The unique name of the policy.
* `cloudlet_code` - (Required) The two- or three- character code for the type of Cloudlet. Enter `ALB` for Application Load Balancer, `AP` for API Prioritization, `AS` for Audience Segmentation, `CD` for Phased Release, `ER` for Edge Redirector, `FR` for Forward Rewrite, `IG` for Request Control, `IV` for Input Validation, or `VP` for Visitor Prioritization.
* `description` - (Optional) The description of this specific policy.
* `group_id` - (Required) Defines the group association for the policy. You must have edit privileges for the group.
* `match_rule_format` - (Optional) The version of the Cloudlet-specific `match_rules`.
* `match_rules` - (Optional) A JSON structure that defines the rules for this policy. See the [Terraform syntax documentation](https://www.terraform.io/docs/configuration-0-11/syntax.html) for more information on embedding multiline strings.

## Attribute reference

The following attributes are returned:

* `cloudlet_id` - A unique identifier that corresponds to a Cloudlets policy type. Enter `0` for Edge Redirector, `1` for Visitor Prioritization, `3` for Forward Rewrite, `4` for Request Control, `5` for API Prioritization, `6` for Audience Segmentation, `7` for Phased Release, `8` for Input Validation, or `9` for Application Load Balancer.
* `version` - The version number of the policy.
* `warnings` - A JSON-encoded list of warnings.

## Import

Basic usage:

```hcl
resource "akamai_cloudlets_policy" "example" {
    # (resource arguments)
  }
```

You can import your Akamai Cloudlets policy using a policy name.

For example:

```shell
$ terraform import akamai_cloudlets_policy.example policy1
```
