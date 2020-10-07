---
layout: "akamai"
page_title: "Akamai: RatePolicies"
subcategory: "APPSEC"
description: |-
 RatePolicies
---

# akamai_appsec_rate_policies

Use `akamai_appsec_rate_policies` data source to retrieve a rate_policies id.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}
data "akamai_appsec_configuration" "appsecconfigedge" {
  name = "Example for EDGE"
  
}



output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}


data "akamai_appsec_rate_policies" "appsecreatepolicies" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}


```

## Argument Reference

The following arguments are supported:

* `config_id`- (Required) The Configuration ID

* `version` - (Required) The Version Number of configuration

# Attributes Reference

The following are the return attributes:

*`output_text` - The rate policies in  formatted text

