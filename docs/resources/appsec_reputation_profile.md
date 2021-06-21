---
layout: "akamai"
page_title: "Akamai: Reputation Profile"
subcategory: "Application Security"
description: |-
 Reputation Profile
---

# akamai_appsec_reputation_profile

Use the `akamai_appsec_reputation_profile` resource to create or modify a reputation profile for a specific security configuration.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

// USE CASE: user wants to create a reputation profile for a given configuration and version, using a JSON definition
resource "akamai_appsec_reputation_profile" "reputation_profile" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  reputation_profile =  file("${path.module}/reputation_profile.json")
}
output "reputation_profile_id" {
  value = akamai_appsec_reputation_profile.reputation_profile_id
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `reputation_profile` - (Required) The name of a file containing a JSON-formatted definition of the reputation profile. ([format](https://developer.akamai.com/api/cloud_security/application_security/v1.html#postreputationprofiles))


## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `reputation_profile_id` - The ID of the newly created or modified reputation profile.

