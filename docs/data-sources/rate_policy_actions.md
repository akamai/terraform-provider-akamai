---
layout: "akamai"
page_title: "Akamai: RatePolicyActions"
subcategory: "APPSEC"
description: |-
 RatePolicyActions
---

# akamai_appsec_rate_policy_actions

Use `akamai_appsec_rate_policy_actions` data source to retrieve a rate_policy_actions id.

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

data "akamai_appsec_rate_policy_actions" "appsecreatepolicysactions" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
    
}

output "ds_rate_policy_actions" {
  value = data.akamai_appsec_rate_policy_actions.appsecreatepolicysactions.output_text
}


```

## Argument Reference

The following arguments are supported:

* `config_id`- (Required) The Configuration ID

* `versionnumber` - (Required) The Version Number of configuration

# Attributes Reference

The following are the return attributes:

*`targetid` - The TargetID

*`output_text` - The rate policies in  formatted text

