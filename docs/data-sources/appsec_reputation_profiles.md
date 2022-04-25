---
layout: "akamai"
page_title: "Akamai: Reputation Profiles"
subcategory: "Application Security"
description: |-
 Reputation Profiles
---

# akamai_appsec_reputation_profiles

**Scopes**: Security configuration; reputation profile

Returns information about your reputation profiles. Reputation profiles grade the security risk of an IP address based on previous activities associated with that address. Depending on the reputation score, and depending on how your configuration has been set up, requests from a specific IP address can trigger an alert, or even be blocked.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/reputation-profiles](https://techdocs.akamai.com/application-security/reference/get-reputation-profiles)

## Example Usage

Basic usage:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}
data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

// USE CASE: User wants to view all the reputation profiles associated with a security configuration.

data "akamai_appsec_reputation_profiles" "reputation_profiles" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}
output "reputation_profiles_output" {
  value = data.akamai_appsec_reputation_profiles.reputation_profiles.output_text
}
output "reputation_profiles_json" {
  value = data.akamai_appsec_reputation_profiles.reputation_profiles.json
}

// USE CASE: User wants to view a specific reputation profile associated with a given configuration

data "akamai_appsec_reputation_profiles" "reputation_profile" {
  config_id             = data.akamai_appsec_configuration.configuration.config_id
  reputation_profile_id = "12345"
}
output "reputation_profile_json" {
  value = data.akamai_appsec_reputation_profiles.reputation_profile.json
}
output "reputation_profile_output" {
  value = data.akamai_appsec_reputation_profiles.reputation_profile.output_text
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the reputation profiles.
- `reputation_profile_id` (Optional). Unique identifier of the reputation profile you want to return information for. If not included, information is returned for all your reputation profiles.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report of the details about the specified reputation profile or profiles.
- `json`. JSON-formatted report of the details about the specified reputation profile or profiles.