---
layout: "akamai"
page_title: "Akamai: RatePolicy"
subcategory: "APPSEC"
description: |-
  RatePolicy
---

# resource_akamai_appsec_rate_policy


The `resource_akamai_appsec_rate_policy` resource allows you to create or re-use RatePolicys.

If the RatePolicy already exists it will be used instead of creating a new one.

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


resource "akamai_appsec_rate_policy" "appsecratepolicy" {
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

*`rate_policy_id` - The Rate Policy ID

