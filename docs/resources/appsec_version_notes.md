---
layout: "akamai"
page_title: "Akamai: VersionNotes"
subcategory: "Application Security"
description: |-
 VersionNotes
---

# akamai_appsec_version_notes

**Scopes**: Security configuration

Updates the version notes for a security configuration.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/version-notes](https://developer.akamai.com/api/cloud_security/application_security/v1.html#putversionnotes)

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

// USE CASE: User wants to update the version notes for the latest version of a security configuration.

resource "akamai_appsec_version_notes" "version_notes" {
  config_id     = data.akamai_appsec_configuration.configuration.config_id
  version_notes = "This version enables reputation profiles."
}
output "version_notes" {
  value = akamai_appsec_version_notes.version_notes.output_text
}
```

## Argument Reference

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration whose version notes are being modified.
- `version_notes` (Required). Brief description of the security configuration version.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `output_text`. Tabular report showing the updated version notes.

