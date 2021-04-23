---
layout: "akamai"
page_title: "Akamai: SecurityPolicy"
subcategory: "Application Security"
description: |-
 SecurityPolicy
---

# akamai_appsec_security_policy

Use the `akamai_appsec_security_policy` data source to retrieve information about the security policies associated with a specific security configuration, or about a specific security policy.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

data "akamai_appsec_security_policy" "security_policies" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "security_policies_list" {
  value = data.akamai_appsec_security_policy.security_policies.security_policy_id_list
}

output "security_policies_text" {
  value = data.akamai_appsec_security_policy.security_policies.output_text
}

data "akamai_appsec_security_policy" "specific_security_policy" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  security_policy_name = "APIs"
}

output "specific_security_policy_id" {
  value = data.akamai_appsec_security_policy.specific_security_policy.security_policy_id
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `security_policy_name`- (Optional) The name of the security policy to use. If not supplied, information about all security policies is returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `security_policy_id_list` - A list of the IDs of all security policies.

* `output_text` - A tabular display showing the ID and name of all security policies.

* `security_policy_id` - The ID of the security policy. Included only if `security_policy_name` was specified.

