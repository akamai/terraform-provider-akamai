---
layout: "akamai"
page_title: "Akamai: Attack Group Condition & Exception"
subcategory: "Application Security"
description: |-
 Attack Group Condition & Exception
---

# akamai_appsec_attack_group_condition_exception

Use the `akamai_appsec_attack_group_condition_exception` resource to create or modify an attack group's conditions and exceptions.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to add condition-exception to an attack grooup using a JSON input
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_attack_group_condition_exception" "condition_exception" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
  attack_group = var.attack_group
  condition_exception = file("${path.module}/condition_exception.json")
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `attack_group` - The attack group to use.

* `condition_exception` - (Required) The name of a file containing a JSON-formatted description of the conditions and exceptions to use ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putattackgroupconditionexception))


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* None

