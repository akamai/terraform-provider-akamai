---
layout: "akamai"
page_title: "Akamai: CustomDeny"
subcategory: "Application Security"
description: |-
 CustomDeny
---

# akamai_appsec_custom_deny

Use the `akamai_appsec_custom_deny` data source to retrieve information about custom deny actions for a specific security configuration version, or about a particular custom deny action. The information available is described [here](https://developer.akamai.com/api/cloud_security/application_security/v1.html#getcustomdeny).

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: user wants to view the custom deny data with a given security configuration
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}

data "akamai_appsec_custom_deny" "custom_deny_list" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
}

//tabular data with id and name
output "custom_deny_list_output" {
  value = data.akamai_appsec_custom_deny.custom_deny_list.output_text
}

output "custom_deny_list_json" {
  value = data.akamai_appsec_custom_deny.custom_deny_list.json
}

// USE CASE: user wants to see a single custom deny associated with a given security configuration version
data "akamai_appsec_custom_deny" "custom_deny" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  custom_deny_id = var.custom_deny_id
}

output "custom_deny_json" {
  value = data.akamai_appsec_custom_deny.custom_deny.json
}

output "custom_deny_output" {
  value = data.akamai_appsec_custom_deny.custom_deny.output_text
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The configuration ID to use.

* `version` - (Required) The version number of the configuration to use.

* `custom_deny_id` - (Optional) The ID of a specific custom deny action.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `json` - A JSON-formatted list of the custom deny action information. 

* `output_text` - A tabular display showing the custom deny action information.


