---
layout: "akamai"
page_title: "Akamai: ExportConfiguration"
subcategory: "APPSEC"
description: |-
 ExportConfiguration
---

# akamai_appsec_export_configuration

Use `akamai_appsec_export_configuration` data source to retrieve a export_configuration id.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Example for EDGE"
}

data "akamai_appsec_export_configuration" "appsecexportconfiguration" {
   config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
   version  = data.akamai_appsec_configuration.appsecconfigedge.latest_version 
search = ["ruleActions","customRules","rulesets","reputationProfiles","ratePolicies","matchTargets"]
}


```

## Argument Reference

The following arguments are supported:

* `config_id`- (Required) The Configuration ID

* `version` - (Required) The Version Number of configuration

* `search` - (Optional) The Version Number of configuration

# Attributes Reference

The following are the return attributes:

* `json` - Export of Configuration data

* `output_text` - Export of Configuration data in tabular format


