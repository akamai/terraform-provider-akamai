---
layout: "akamai"
page_title: "Akamai: Rate Policy Action"
subcategory: "Application Security"
description: |-
  Rate Policy Action
---

# resource_akamai_appsec_rate_policy_action


(Beta) The `resource_akamai_appsec_rate_policy_action` resource allows you to create, modify or delete the actions in a rate policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to create a rate policy and rate policy actions for a given security configuration and version
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_rate_policy" "appsec_rate_policy" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  rate_policy =  file("${path.module}/rate_policy.json")
}
resource  "akamai_appsec_rate_policy_action" "appsec_rate_policy_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  rate_policy_id = akamai_appsec_rate_policy.appsec_rate_policy.rate_policy_id
  ipv4_action = "deny"
  ipv6_action = "deny"
}

//TF destroy - means set the rate policy action to none
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `rate_policy_id` - (Required) The ID of the rate policy to use.

* `ipv4_action` - (Required) The ipv4 action to assign to this rate policy, either `alert`, `deny`, or `none`. If the action is none, the rate policy is inactive in the policy.

* `ipv6_action` - (Required) The ipv6 action to assign to this rate policy, either `alert`, `deny`, or `none`. If the action is none, the rate policy is inactive in the policy.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* [TBD; currently None] `rate_policy_action_id` - The ID of the modified or newly created rate policy action.

