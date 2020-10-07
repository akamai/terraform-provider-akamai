---
layout: "akamai"
page_title: "Akamai: RatePolicyAction"
subcategory: "APPSEC"
description: |-
  RatePolicyAction
---

# resource_akamai_appsec_rate_policy_action


The `resource_akamai_appsec_rate_policy_action` resource allows you to create or re-use RatePolicyActions.

If the RatePolicyAction already exists it will be used instead of creating a new one.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}
data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Example for EDGE"
  
}



output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}


resource "akamai_appsec_rate_policy_action" "appsecreatepolicyaction" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "f1rQ_106946"
    rate_policy_id = 321456
    ipv4_action = "alert"
    ipv6_action = "none"
}

output "ratepolicyaction" {
  value = akamai_appsec_rate_policy_action.appsecreatepolicyaction.rate_policy_id
}

```

## Argument Reference

The following arguments are supported:
* `config_id`- (Required) The Configuration ID

* `version` - (Required) The Version Number of configuration

* `policy_id` - (Required) The Policy Id of configuration

* `ipv4_action` - (Required) The ipv4action for rate policy action

* `ipv6_action` - (Required) The ipv6action for rate policy action

# Attributes Reference

The following are the return attributes:

*`rate_policy_id` - The Rate Policy ID

