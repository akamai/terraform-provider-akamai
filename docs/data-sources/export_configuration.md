---
layout: "akamai"
page_title: "Akamai: ExportConfiguration"
subcategory: "APPSEC"
description: |-
 ExportConfiguration
---

# akamai_appsec_export_configuration

Use the akamai_appsec_export_configuration` data source to retrieve comprehensive details about a security configuration and version, including rate and security policies, rules, hostnames, and other settings. You can retrieve the entire set of information in JSON format, or a subset of the information in tabular format.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "configuration" {
  name = "Akamai Tools"
}

data "akamai_appsec_export_configuration" "export" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  search = ["securityPolicies", "selectedHosts"]
}

output "json" {
  value = data.akamai_appsec_export_configuration.export.json
}

output "text" {
  value = data.akamai_appsec_export_configuration.export.output_text
}

```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `search` - (Optional) A bracket-delimited list of quoted strings specifying the types of information to be retrieved and made available for display in the `output_text` format. The following types are available:
  * customRules
  * matchTargets
  * ratePolicies
  * reputationProfiles
  * rulesets
  * securityPolicies
  * selectableHosts
  * selectedHosts

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - The complete set of information about the specified security configuration version, in JSON format. This includes the types available for the `search` parameter, plus several additional fields such as createDate and createdBy.

* `output_text` - A tabular display showing the types of data specified in the `search` parameter. Included only if the `search` parameter specifies at least one type.

