---
layout: "akamai"
page_title: "Akamai: Policy Protections"
subcategory: "Application Security"
description: |-
 Policy Protections
---

# akamai_appsec_waf_protections

Use the `akamai_appsec_waf_protections` data source to retrieve the protections in place for a given security configuration version and security policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to view all security policy protections
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_security_policy_protections" "protections" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
}
output "protections_json" {
  value = data.akamai_appsec_security_policy_protections.protections.json
}
output "protections_applyApiConstraints" {
  value = data.akamai_appsec_security_policy_protections.protections.apply_api_constraints
}
output "protections_applyApplicationLayerControls" {
  value = data.akamai_appsec_security_policy_protections.protections.apply_application_layer_controls
}
output "protections_applyBotmanControls" {
  value = data.akamai_appsec_security_policy_protections.protections.apply_botman_controls
}
output "protections_applyNetworkLayerControls" {
  value = data.akamai_appsec_security_policy_protections.protections.apply_network_layer_controls
}
output "protections_applyRateControls" {
  value = data.akamai_appsec_security_policy_protections.protections.apply_rate_controls
}
output "protections_applyReputationControls" {
  value = data.akamai_appsec_security_policy_protections.protections.apply_reputation_controls
}
output "protections_applySlowPostControls" {
  value = data.akamai_appsec_security_policy_protections.protections.apply_slow_post_controls
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - a JSON-formatted list of the protections

* `apply_api_constraints` - true if api constraints are applied, otherwise false 

* `apply_application_layer_controls` - true if application layer controls are applied, otherwise false 

* `apply_botman_controls` - true if botman controls are applied, otherwise false 

* `apply_network_layer_controls` - true if network layer controls are applied, otherwise false 

* `apply_rate_controls` - true if rate controls are applied, otherwise false 

* `apply_reputation_controls` - true if reputation controls are applied, otherwise false 

* `apply_slow_post_controls` - true if slow post controls are applied, otherwise false 
