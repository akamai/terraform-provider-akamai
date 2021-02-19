---
layout: "akamai"
page_title: "Akamai: Security Policy Protections"
subcategory: "Application Security"
description: |-
 Security Policy Protections
---

# akamai_appsec_security_policy_protections

Use the `akamai_appsec_security_policy_protections` data source to retrieve the protections in effect for a given security policy.

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

* `apply_application_layer_controls` - `true` or `false`, indicating whether application layer controls are in effect.

* `apply_network_layer_controls` - `true` or `false`, indicating whether network layer controls are in effect.

* `apply_rate_controls` - `true` or `false`, indicating whether rate controls are in effect.

* `apply_reputation_controls` - `true` or `false`, indicating whether reputation controls are in effect.

* `apply_botman_controls` - `true` or `false`, indicating whether botman controls are in effect.

* `apply_api_constraints` - `true` or `false`, indicating whether API constraints are in effect.

* `apply_slow_post_controls` - `true` or `false`, indicating whether slow post controls are in effect.

* `json` - a JSON-formatted list showing the status of the protection settings

* `output_text` - a tabular display showing the status of the protection settings

