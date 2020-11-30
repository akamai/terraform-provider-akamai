---
layout: "akamai"
page_title: "Akamai: FOO"
subcategory: "Application Security"
description: |-
 FOO
---

# akamai_appsec_FOO

Use the `akamai_appsec_FOO` resource to create or modify ...

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}


```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

* `foo` - (Required) The name of a file containing a JSON-formatted ([format]())

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the ID, name, and action of all custom rules associated with the specified security configuration, version and security policy.


