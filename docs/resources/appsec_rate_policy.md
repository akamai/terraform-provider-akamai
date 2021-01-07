---
layout: "akamai"
page_title: "Akamai: Rate Policy"
subcategory: "Application Security"
description: |-
  Rate Policy
---

# resource_akamai_appsec_rate_policy


The `resource_akamai_appsec_rate_policy` resource allows you to create, modify or delete rate policies for a specific security configuration version.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to create a rate policy for a given configuration and version, using a JSON rule definition
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_rate_policy" "rate_policy" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  rate_policy =  file("${path.module}/rate_policy.json")
}
output "rate_policy_id" {
  value = akamai_appsec_rate_policy.rate_policy.rate_policy_id
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `rate_policy` - (Required) The name of a file containing a JSON-formatted rate policy definition ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#57c65cbd)).

* `rate_policy_id` - (Optional) The ID of an existing rate policy to be modified.


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `rate_policy_id` - The ID of the modified or newly created rate policy.

