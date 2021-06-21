---
layout: "akamai"
page_title: "Akamai: SecurityPolicyRename"
subcategory: "Application Security"
description: |-
  SecurityPolicyRename
---

# akamai_appsec_security_policy_rename

The `akamai_appsec_security_policy_rename` resource allows you to rename an existing security policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to rename a security policy
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_security_policy" "security_policy_rename" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = var.security_policy_id
  security_policy_name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to be renamed.

* `security_policy_name` - (Required) The new name to be given to the security policy.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* None

