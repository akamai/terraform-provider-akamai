---
layout: "akamai"
page_title: "Akamai: WAF Protection"
subcategory: "Application Security"
description: |-
 WAF Protection
---

# akamai_appsec_waf_protection

Use the `akamai_appsec_waf_protection` data source to retrieve the current protection settings for a given security configuration version and policy


## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_waf_protection" "waf_protection" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  policy_id = var.policy_id
}
output "output_text" {
  value = data.akamai_appsec_waf_protection.waf_protection.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `policy_id` - (Required) The ID of the security policy to use

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the enabled status (true or false) of the following protection features:
  * applyApiConstraints
  * applyApplicationLayerControls
  * applyBotmanControls
  * applyNetworkLayerControls
  * applyRateControls
  * applyReputationControls
  * applySlowPostControls
	
