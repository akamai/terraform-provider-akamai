---
layout: "akamai"
page_title: "Akamai: Security Policy Protections"
subcategory: "Application Security"
description: |-
 Security Policy Protections
---

# akamai_appsec_security_policy_protections

Use the `akamai_appsec_security_policy_protections` resource to create or modify ...

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to update the security policy protection settings
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

resource "akamai_appsec_security_policy_protections" "protections" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  apply_application_layer_controls = var.apply_application_layer_controls
  apply_network_layer_controls = var.apply_network_layer_controls
  apply_rate_controls = var.apply_rate_controls
  apply_reputation_controls = var.apply_reputation_controls
  apply_botman_controls = var.apply_botman_controls
  apply_api_constraints = var.apply_api_constraints
  apply_slow_post_controls = var.apply_slow_post_controls
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `apply_application_layer_controls` - (Required) Whether to enable application layer controls: either `true` or `false`.

* `apply_network_layer_controls` - (Required) Whether to enable network layer controls: either `true` or `false`.

* `apply_rate_controls` - (Required) Whether to enable rate controls: either `true` or `false`.

* `apply_reputation_controls` - (Required) Whether to enable reputation controls: either `true` or `false`.

* `apply_botman_controls` - (Required) Whether to enable botman controls: either `true` or `false`.

* `apply_api_constraints` - (Required) Whether to enable api constraints: either `true` or `false`.

* `apply_slow_post_controls` - (Required) Whether to enable slow post controls: either `true` or `false`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the protection settings in effect.

