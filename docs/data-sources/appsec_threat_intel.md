---
layout: "akamai"
page_title: "Akamai: Threat Intelligence
subcategory: "Application Security"
description: |-
 Threat Intelligence
---

# akamai_appsec_threat_intel

Use the `akamai_appsec_threat_intel` data source to view threat intelligence setting for a policy
__BETA__ This is Adaptive Security Engine(ASE) related data source. Please contact your akamai representative if you want to learn more

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view threat intelligence setting for a policy
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_threat_intel" "threat_intel" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
}
output "threat_intel" {
  value = data.akamai_appsec_threat_intel.threat_intel.threat_intel
}

output "json" {
  value = data.akamai_appsec_threat_intel.threat_intel.json
}
output "output_text" {
  value = data.akamai_appsec_threat_intel.threat_intel.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `threat_intel` - Threat Intelligence setting, either `on` or `off`.

* `json` - A JSON-formatted threat intelligence object

* `output_text` - A tabular display of the threat intel information.



