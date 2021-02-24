---
layout: "akamai"
page_title: "Akamai: Reputation Profile Actions"
subcategory: "Application Security"
description: |-
 Reputation Profile Actions
---

# akamai_appsec_reputation_profile_actions

Use the `akamai_appsec_reputation_profile_actions` data source to retrieve details about reputation profiles and their associated actions, or about the actions associated with a specific reputation profile.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the reputation profile actions associated with a given security policy
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
data "akamai_appsec_reputation_profile_actions" "reputation_profile_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
}
output "reputation_profile_actions_text" {
  value = data.akamai_appsec_reputation_profile_actions.reputation_profile_actions.output_text
}
output "reputation_profile_actions_json" {
  value = data.akamai_appsec_reputation_profile_actions.reputation_profile_actions.json
}

// USE CASE: user wants to view the action for a single reputation profile associated with a given security policy
data "akamai_appsec_reputation_profile_actions" "reputation_profile_actions" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  reputation_profile_id = var.reputation_profile_id
}

output "reputation_profile_action" {
  value = data.akamai_appsec_reputation_profile_actions.reputation_profile_actions.action
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) THe ID of the security policy to use.

* `reputation_profile_id` - (Optional) The ID of a given reputation profile. If not supplied, information about all reputation profiles is returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `action` - The action that the specified reputation profile or profiles take when triggered.

* `json` - A JSON-formatted display of the specified reputation profile action information.

* `output_text` - A tabular display of the specified reputation profile action information.

