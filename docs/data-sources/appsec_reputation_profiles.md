---
layout: "akamai"
page_title: "Akamai: Reputation Profiles"
subcategory: "Application Security"
description: |-
 Reputation Profiles
---

# akamai_appsec_reputation_profiles

Use the `akamai_appsec_reputation_profiles` data source to retrieve details about all reputation profiles, or a specific reputation profiles.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

// USE CASE: user wants to see all the reputation profiles associated with a given configuration and version, or a single reputation profile.
data "akamai_appsec_reputation_profiles" "reputation_profiles" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
}
output "reputation_profiles_output" {
  value = data.akamai_appsec_reputation_profiles.reputation_profiles.output_text
}
output "reputation_profiles_json" {
  value = data.akamai_appsec_reputation_profiles.reputation_profiles.json
}

// USE CASE: user wants to see a single reputation profile associated with a given configuration and version
data "akamai_appsec_reputation_profiles" "reputation_profile" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  reputation_profile_id = var.reputation_profile_id
}
output "reputation_profile_json" {
  value = data.akamai_appsec_reputation_profiles.reputation_profile.json
}
output "reputation_profile_output" {
  value = data.akamai_appsec_reputation_profiles.reputation_profile.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `reputation_profile_id` - (Optional) The ID of a given reputation profile. If not supplied, information about all reputation profiles is returned.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display of the details about the indicated reputation profile or profiles.

* `json` - A JSON-formatted display of the details about the indicated reputation profile or profiles.

