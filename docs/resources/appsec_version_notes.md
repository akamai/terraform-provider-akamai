---
layout: "akamai"
page_title: "Akamai: VersionNotes"
subcategory: "Application Security"
description: |-
 VersionNotes
---

# akamai_appsec_version_notes

Use the `akamai_appsec_version_notes` resource to update the version notes for a configuration.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

// USE CASE: user wants to update the version notes of the latest version
resource "akamai_appsec_version_notes" "version_notes" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version_notes = var.version_notes
}
output "version_notes" {
  value = akamai_appsec_version_notes.version_notes.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The configuration ID to use.

* `version_notes` - (Required) A string containing the version notes to be used.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `output_text` - A tabular display showing the updated version notes.

