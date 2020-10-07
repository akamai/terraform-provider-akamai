---
layout: "akamai"
page_title: "Akamai: CustomRuleAction"
subcategory: "APPSEC"
description: |-
  CustomRuleAction
---

# resource_akamai_appsec_custom_rule_action


The `resource_akamai_appsec_custom_rule_action` resource allows you to create or re-use CustomRuleActions.

If the CustomRuleAction already exists it will be used instead of creating a new one.

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


resource "akamai_appsec_custom_rule_action" "appsecreatecustomruleaction" {
    config_id = data.akamai_appsec_configuration.appsecconfigedge.config_id
    version = data.akamai_appsec_configuration.appsecconfigedge.latest_version
    policy_id = "AAAA_81230"
    rule_id = akamai_appsec_custom_rule.appseccustomrule1.rule_id
    custom_rule_action = "alert"
}

output "customruleaction" {
  value = akamai_appsec_custom_rule_action.appsecreatecustomruleaction.rule_id
}

```

## Argument Reference

The following arguments are supported:
* `config_id`- (Required) The Configuration ID

* `version` - (Required) The Version Number of configuration

* `policy_id` - (Required) The Policy Id of configuration

* `rule_id` - (Required) The Rule Id of configuration

* `custom_rule_action` - (Required) The custom_rule_action for custom rules  action

# Attributes Reference

The following are the return attributes:

*`rule_id` - The Rule ID

