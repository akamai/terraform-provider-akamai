---
layout: "akamai"
page_title: "Akamai: Rule Upgrade"
subcategory: "Application Security"
description: |-
 Rule Upgrade
---

TBD
# akamai_appsec_rule_upgrade

Use the `akamai_appsec_rule_upgrade` resource to upgrade to the most recent version of the KRS rule set. Akamai periodically updates these rules to keep protections current. However, the rules you use in your security policies do not automatically upgrade to the latest version when using mode: KRS. These rules do update automatically when you have mode set to AAG. Before you upgrade, run Get upgrade details to see which rules have changed. If you want to test how these rules would operate with live traffic before committing to the upgrade, run them in evaluation mode. This applies to KRS rules only and does not allow you to make any changes to the rules themselves. The response is the same as the mode response. 

## Example Usage

Basic usage:

```hcl
provider "akamai" {
  appsec_section = "default"
}

// USE CASE: user wants to set the waf mode
data "akamai_appsec_configuration" "configuration" {
  name = var.security_configuration
}
resource "akamai_appsec_rule_upgrade" "rule_upgrade" {
  config_id = data.akamai_appsec_configuration.configuration.config_id
  version = data.akamai_appsec_configuration.configuration.latest_version
  security_policy_id = var.security_policy_id
}
output "rule_upgrade_current_ruleset" {
  value = akamai_appsec_rule_upgrade.rule_upgrade.current_ruleset
}
output "rule_upgrade_mode" {
  value = akamai_appsec_rule_upgrade.rule_upgrade.mode
}
output "rule_upgrade_eval_status" {
  value = akamai_appsec_rule_upgrade.rule_upgrade.eval_status
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The ID of the security configuration to use.

* `version` - (Required) The version number of the security configuration to use.

* `security_policy_id` - (Required) The ID of the security policy to use.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

 * `current_ruleset` - A string indicating the version number and release date of the current KRS rule set.

 * `mode` - A string indicating the current mode, either "KRS" or "AAG".

 * `eval_status` - TBD

