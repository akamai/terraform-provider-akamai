---
layout: "akamai"
page_title: "Akamai: Reputation Profile Action"
subcategory: "Application Security"
description: |-
 Reputation Profile Action
---

# akamai_appsec_reputation_profile_action

(Beta) Use the `akamai_appsec_reputation_profile_action` resource to update what action should be taken when a reputation profile's rule is triggered.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// OPEN API --> https://developer.akamai.com/api/cloud_security/application_security/v1.html#putreputationprofileaction

data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

resource  "akamai_appsec_reputation_profile_action" "appsec_reputation_profile_action" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  reputation_profile_id = akamai_appsec_reputation_profile.reputation_profile.id
  action = "alert"
}

output "reputation_profile_id" {
  value = akamai_appsec_reputation_profile.reputation_profile.reputation_profile_id
}

output "reputation_profile_action" {
  value = akamai_appsec_reputation_profile_action.appsec_reputation_profile_action.action
}

//TF destroy - means set the reputation profile action to none

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `reputation_profile_id` - (Required) The ID of the reputation profile to use.

* `action` - (Required) The action to take when the specified reputation profile’s rule is triggered: `alert` to record the trigger event, `deny` to block the request, or `none` to take no action.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

