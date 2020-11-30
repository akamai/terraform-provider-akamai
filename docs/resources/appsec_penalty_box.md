---
layout: "akamai"
page_title: "Akamai: Penalty Box"
subcategory: "Application Security"
description: |-
 Penalty Box
---

# akamai_appsec_penalty_box

Use the `akamai_appsec_penalty_box` resource to update the penalty box settings for a given security policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// OPEN API --> https://developer.akamai.com/api/cloud_security/application_security/v1.html#putpenaltybox

// USE CASE: user wants to set the penalty box
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_penalty_box" "penalty_box" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  penalty_box_protection = true
  penalty_box_action = var.action
}

//TF destroy - set the penalty_box_protection to false.

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `penalty_box_protection` - (Required) A boolean value indicating whether to enable penalty box protection.

* `penalty_box_action` - (Required) "Deny", "Alert" or "None", indicating the action to take when penalty box protection is triggered. Ignored if `penalty_box_protection` is set to `false`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the updated penalty box protection settings.


