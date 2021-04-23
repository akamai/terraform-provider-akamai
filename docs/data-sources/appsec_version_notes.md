---
layout: "akamai"
page_title: "Akamai: VersionNotes"
subcategory: "Application Security"
description: |-
 VersionNotes
---

# akamai_appsec_version_notes

Use the `akamai_appsec_version_notes` data source to retrieve the most recent version notes for a configuration.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

// USE CASE: user wants to see version notes of the latest version
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

The following arguments are supported:

* `config_id` - (Required) The configuration ID to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list showing the version notes.

* `output_text` - A tabular display showing the version notes.

