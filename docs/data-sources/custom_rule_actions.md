---
layout: "akamai"
page_title: "Akamai: CustomRuleActions"
subcategory: "APPSEC"
description: |-
 CustomRuleActions
---

# akamai_appsec_custom_rule_actions

Use `akamai_appsec_custom_rule_actions` data source to retrieve a custom_rule_actions id.

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}


data "akamai_appsec_custom_rule_actions" "appsecreatecustomruleactions" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
}

output "customruleactions" {
  value = data.akamai_appsec_custom_rule_actions.appsecreatecustomruleactions.output_text
}


```

## Argument Reference

The following arguments are supported:

* `config_id`- (Required) The Configuration ID

* `version` - (Required) The Version Number of configuration

* `policy_id` - (Required) The Policy Id of configuration

# Attributes Reference

The following are the return attributes:

*`rule_id` - The Rule ID

*`ouput_text` - The list of custom rule actions in tabular format

