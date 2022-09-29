---
layout: akamai
subcategory: Application Security
---

# akamai_appsec_version_notes

**Scopes**: Security configuration

Returns the most recent version notes for a security configuration.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/version-notes](https://techdocs.akamai.com/application-security/reference/get-version-notes)

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

// USE CASE: User wants to view the version notes for the most-recent version of a security configuration.

data "akamai_appsec_version_notes" "version_notes" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
}

output "version_notes_text" {
  value = data.akamai_appsec_version_notes.version_notes.output_text
}

output "version_notes_json" {
  value = data.akamai_appsec_version_notes.version_notes.json
}
```

## Argument Reference

This data source supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration you want to return information for.

## Output Options

The following options can be used to determine the information returned, and how that returned information is formatted:

- `json`. JSON-formatted list showing the version notes.
- `output_text`. Tabular report showing the version notes.