---
layout: "akamai"
page_title: "Akamai: Reputation Profile"
subcategory: "Application Security"
description: |-
 Reputation Profile
---

# akamai_appsec_reputation_profile

**Scopes**: Security policy

Creates or modifies a reputation profile.
Reputation profiles grade the security risk of an IP address based on previous activities associated with that address.
Depending on the reputation score and how your configuration has been set up, requests from a specific IP address can trigger an alert or even be blocked.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/reputation-profiles](https://techdocs.akamai.com/application-security/reference/put-reputation-profile)

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

// USE CASE: User wants to create a reputation profile for a given security configuration by using a JSON-formatted definition.

resource "akamai_appsec_reputation_profile" "reputation_profile" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  reputation_profile = file("${path.module}/reputation_profile.json")
}
output "reputation_profile_id" {
  value = akamai_appsec_reputation_profile.reputation_profile.reputation_profile_id
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the reputation profile being modified.
- `reputation_profile` (Required). Path to a JSON file containing a definition of the reputation profile. You can view a sample JSON file in the [Create a reputation profile](https://developer.akamai.com/api/cloud_security/application_security/v1.html#postreputationprofiles) section of the Application Security API documentation.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `reputation_profile_id`. ID of the newly-created or newly-modified reputation profile.