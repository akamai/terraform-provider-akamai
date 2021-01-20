---
layout: "akamai"
page_title: "Akamai: ReputationAnalysis"
subcategory: "Application Security"
description: |-
 ReputationAnalysis
---

# akamai_appsec_reputation_analysis

Use the `` data source to retrieve information about the current reputation analysis settings. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getreputationanalysis).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the all reputation analysis associated with a given security policy
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_reputation_analysis" "reputation_analysis" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
}

output "reputation_analysis_text" {
  value = data.akamai_appsec_reputation_analysis.reputation_analysis.output_text
}

output "reputation_analysis_json" {
  value = data.akamai_appsec_reputation_analysis.reputation_analysis.json
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The configuration ID to use.

* `version` - (Required) The version number of the configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list of the reputation analysis settings.

* `output_text` - A tabular display showing the reputation analysis settings.

