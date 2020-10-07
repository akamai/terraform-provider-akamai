---
layout: "akamai"
page_title: "Akamai: SecurityPolicy"
subcategory: "APPSEC"
description: |-
 SecurityPolicy
---

# akamai_appsec_security_policy

Use `akamai_appsec_security_policy` data source to retrieve a security_policy id.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}


data "akamai_appsec_configuration_version" "appsecconfigurationversion" {
    name = "Akamai Tools"
   }

output "configsedge" {
  value = data.akamai_appsec_configuration.appsecconfigedge.config_id
}

output "configsedgelatestversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.latest_version
}

output "configsedgeconfiglist" {
  value = data.akamai_appsec_configuration.appsecconfigedge.output_text
}

output "configsedgeconfigversion" {
  value = data.akamai_appsec_configuration.appsecconfigedge.version
}
data "akamai_appsec_security_policy" "appsecsecuritypolicy" {
  name = "akamaitools" 
  config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
  version =  data.akamai_appsec_configuration.appsecconfigedge.version
}

output "securitypolicy" {
  value = data.akamai_appsec_security_policy.appsecsecuritypolicy.policy_id
}

output "securitypolicies" {
  value = data.akamai_appsec_security_policy.appsecsecuritypolicy.policy_list
}

```

## Argument Reference

The following arguments are supported:

* `name`- (Optional) The Configuration Name

* `config_id` - (Required) The ID Number of configuration

* `version` - (Required) The ID Number of configuration

# Attributes Reference

The following are the return attributes:

* `policy_id` - Policy Id of configuration

* `policy_list` - Policy list of configuration list

