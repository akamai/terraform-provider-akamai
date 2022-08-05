---
layout: "akamai"
page_title: "Akamai: Security Policy Protections"
subcategory: "Application Security"
description: |-
 Security Policy Protections
---

# akamai_appsec_security_policy_protections

**Scopes**: Security policy

Returns information about the protections in effect for the specified security policy.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/protections](https://techdocs.akamai.com/application-security/reference/get-policy-protections)

## Example Usage

Basic usage:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: User wants to view all security policy protections.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
data "akamai_appsec_security_policy_protections" "protections" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
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

output "protections_applyMalwareControls" {
  value = data.akamai_appsec_security_policy_protections.protections.apply_malware_controls
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

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the security policy protections.
- `security_policy_id` (Required). Unique identifier of the security policy you want to return protections information for.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `apply_application_layer_controls`. Returns **true** if application layer controls are enabled; returns **false** if they are not.
- `apply_api_constraints`. Returns **true** if API constraints are enabled; returns **false** if they are not.
- `apply_botman_controls`. Returns **true** if Bot Manager controls are enabled; returns **false** if they are not.
- `apply_malware_controls`. Returns **true** if malware controls are enabled; returns **false** if they are not.
- `apply_network_layer_controls`. Returns **true** if network layer controls are enabled; returns **false** if they are not.
- `apply_rate_controls`. Returns **true** if rate controls are enabled; returns **false** if they are not.
- `apply_reputation_controls`. Returns **true** if reputation controls are enabled; returns **false** if they are not.
- `apply_slow_post_controls`. Returns **true** if slow POST controls are enabled; returns **false** if they are not.
- `json`. JSON-formatted list showing the status of the protection settings.
- `output_text`. Tabular report showing the status of the protection settings.
