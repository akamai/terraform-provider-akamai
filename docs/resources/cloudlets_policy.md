---
layout: "akamai"
page_title: "Akamai: cloudlets_policy"
subcategory: "Cloudlets"
description: |-
  Cloudlets Policy
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
    "useIncomingSchemeAndHost": true
  },
  {
    "name": "rule2",
    "type": "erMatchRule",
    "matches": [
      {
        "matchType": "hostname",
        "matchValue": "3333.dom",
        "matchOperator": "equals",
        "caseSensitive": true,
        "negate": false
      }
    ],
    "useRelativeUrl": "none",
    "statusCode": 301,
    "redirectURL": "https://www.example.com",
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
* `cloudlet_code` - (Required) The two- or three- character code for the type of Cloudlet, either `ALB` for Application Load Balancer or `ER` for Edge Redirector.
* `description` - (Optional) The description of this specific policy.
* `group_id` - (Required) Defines the group association for the policy. You must have edit privileges for the group.
* `match_rule_format` - (Optional) The version of the Cloudlet-specific `match_rules`.
* `match_rules` - (Optional) A JSON structure that defines the rules for this policy. See the [Terrfaform syntax documentation](https://www.terraform.io/docs/configuration-0-11/syntax.html) for more information on embedding multiline strings.

## Attribute reference

The following attributes are returned:

* `cloudlet_id` - A unique identifier that corresponds to a Cloudlets policy type, either `0` for Edge Redirector or `9` for Application Load Balancer.
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
