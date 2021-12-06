---
layout: "akamai"
page_title: "Akamai: akamai_cloudlets_application_load_balancer_activation"
subcategory: "Cloudlets"
description: |-
  Application Load Balancer activation
---

# akamai_cloudlets_application_load_balancer_activation

Use the `akamai_cloudlets_application_load_balancer_activation` resource to activate the Application Load Balancer Cloudlet configuration. An activation deploys the configuration version to either the Akamai staging or production network. You can activate a specific version multiple times if you need to.

Before activating on production, activate on staging first. This way you can detect any problems in staging before your changes progress to production.

## Example usage

Basic usage:

```hcl
resource "akamai_cloudlets_application_load_balancer_activation" "example" {
  origin_id = "alb_test_1"
  network = "staging"
  version = 1
}
output "status" {
  value = akamai_cloudlets_application_load_balancer_activation.example.status
}
```

## Argument reference

The following arguments are supported:

* `origin_id` - (Required) The identifier of an origin that represents the data center. The Conditional Origin, which is defined in Property Manager, must have an origin type of either `CUSTOMER` or `NET_STORAGE` set in the `origin` behavior. See [property rules](../data-sources/property-rules.md) for more information.
* `network` - (Required) The network you want to activate the policy version on, either `staging`, `stag`,  and `s` for the Staging network, or `production`, `prod`, and `p` for the Production network. All values are case insensitive.
* `version` - (Required) The Application Load Balancer Cloudlet configuration version you want to activate.

## Attribute reference

The following attributes are returned:

* `status` - The activation status for this load balancing configuration.
