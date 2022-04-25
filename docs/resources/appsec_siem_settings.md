---
layout: "akamai"
page_title: "Akamai: SIEMSettings"
subcategory: "Application Security"
description: |-
 SIEMSettings
---

# akamai_appsec_siem_settings

**Scopes**: Security configuration

Modifies SIEM (Security Information and Event Management) integration settings for a security configuration.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/siem](https://techdocs.akamai.com/application-security/reference/put-siem)

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

// USE CASE: User wants to update the SIEM settings.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}

data "akamai_appsec_siem_definitions" "siem_definition" {
  siem_definition_name = "SIEM Version 01"
}

data "akamai_appsec_security_policy" "security_policies" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

resource "akamai_appsec_siem_settings" "siem" {
  config_id               = data.akamai_appsec_configuration.configuration.config_id
  enable_siem             = true
  enable_for_all_policies = false
  enable_botman_siem      = true
  siem_id                 = data.akamai_appsec_siem_definitions.siem_definition.id
  security_policy_ids     = data.akamai_appsec_security_policy.security_policies.security_policy_id_list
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the SIEM settings being modified.
- `enable_siem` (Required). Set to **true** to enable SIEM; set to **false** to disable SIEM.
- `enable_for_all_policies` (Required). Set to **true** to enable SIEM on all security policies in the security configuration; set to **false** to only enable SIEM on the security policies specified by the `security_policy_ids` argument.
- `enable_botman_siem` (Required). Set to **true** to include Bot Manager events in your SIEM events; set to **false** to exclude Bot Manager events from your SIEM events.
- `siem_id` (Required). Unique identifier of the SIEM settings being modified.
- `security_policy_ids` (Required if `enable_for_all_policies` is **false**) JSON array of IDs for the security policies where SIEM integration is to be enabled.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the updated SIEM integration settings.