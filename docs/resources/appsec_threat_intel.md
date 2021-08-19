---
layout: "akamai"
page_title: "Akamai: Threat Intelligence"
subcategory: "Application Security"
description: |-
 Threat Intelligence
---

# akamai_appsec_threat_intel

Use `akamai_appsec_threat_intel` resource to update threat intelligence setting for a policy. Only applies to ASE Manual rulesets. Allowed values are on and off
__BETA__ This is Adaptive Security Engine(ASE) related data resource. Please contact your akamai representative if you want to learn more

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to update threat intelligence setting
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_threat_intel" "threat_intel" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  threat_intel = var.threat_intel
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `threat_intel` - (Required) threat_intel - "on" or "off"  

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None
