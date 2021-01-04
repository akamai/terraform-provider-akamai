---
layout: "akamai"
page_title: "Akamai: Reputation Protection"
subcategory: "Application Security"
description: |-
 Reputation Protection
---

# akamai_appsec_reputation_protection

Use the `akamai_appsec_reputation_protection` resource to enable or disable reputation protection for a given configuration version and security policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to enable or disable reputation protection
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

resource "akamai_appsec_reputation_protection" "protection" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  enabled = var.enabled
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `enabled` - (Required) Whether to enable reputation controls: either `true` or `false`.


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the current protection settings.

