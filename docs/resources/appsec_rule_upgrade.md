---
layout: "akamai"
page_title: "Akamai: Rule Upgrade"
subcategory: "Application Security"
description: |-
 Rule Upgrade
---

# akamai_appsec_rule_upgrade

**Scopes**: Security policy

Upgrades your Kona Rule Set (KRS) rules to the most recent version.
Akamai periodically updates these rules to keep protections current.
However, the rules you use in your security policies are not automatically upgraded to the latest version if you are running in **KRS** or **ASE_MANUAL** mode.
(These rules *do* update automatically when you have mode set to **AAG** or **ASE_AUTO**.)
This resource upgrades your Kona Rule Set rules for organizations running in **KRS** or **ASE_MANUAL** mode.

Note that **ASE_MANUAL **and **ASE_AUTO** modes only apply to organizations running the beta version of Adaptive Security Engine (ASE). Please contact your Akamai representative if you'd like more information about the ASE beta.

Before you upgrade it's recommended that you use the [akamai_appsec_rule_upgrade_details](https://registry.terraform.io/providers/akamai/akamai/latest/docs/data-sources/appsec_rule_upgrade_details) data source to determine which rules and rule sets (if any) have available upgrades. In addition to that, you might want to test the new rules in [evaluation mode](https://registry.terraform.io/providers/akamai/akamai/latest/docs/resources/appsec_eval). In evaluation mode, rules are triggered the same way they are on the production network; however, the only action taken by the rules is to record how they *would* have responded had they been active on the production network. This enables you to see how the rules interact with your production network without actually making changes to that network.

**Related API Endpoint**: [/appsec/v1/configs/{configId}/versions/{versionNumber}/security-policies/{policyId}/rules](https://techdocs.akamai.com/application-security/reference/put-policy-rules)

## Example Usage

Basic usage:

```
terraform {
  required_providers {
    akamai = {
      source = "akamai/akamai"
    }
  }
}

provider "akamai" {
  edgerc = "~/.edgerc"
}

// USE CASE: User wants to set the WAF mode.

data "akamai_appsec_configuration" "configuration" {
  name = "Documentation"
}
resource "akamai_appsec_rule_upgrade" "rule_upgrade" {
  config_id          = data.akamai_appsec_configuration.configuration.config_id
  security_policy_id = "gms1_134637"
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

This resource supports the following arguments:

- `config_id` (Required). Unique identifier of the security configuration associated with the ruleset being upgraded.
- `security_policy_id` (Required). Unique identifier of the security policy associated with the ruleset being upgraded.
- `upgrade_mode`. (Optional). Modifies the upgrade type for organizations running the ASE beta. Allowed values are:
  - **ASE_AUTO**. Akamai automatically updates your rulesets.
  - **ASE_MANUAL**. Manually updates your rulesets.

## Output Options

The following options can be used to determine the information returned and how that returned information is formatted:

- `current_ruleset`. Versioning information for your current KRS rule set.
- `mode`. Specifies the current upgrade mode type. Valid values are:
  - **KRS**. Rulesets must be manually upgraded.

  - **AAG**. Rulesets are automatically upgraded by Akamai.

  - **ASE_MANUAL**. Adaptive Security Engine rulesets must be manually upgraded.

  - **ASE_AUTO**. Adaptive Security Engine rulesets are automatically updated by Akamai.

- `eval_status`. Returns **enabled** if an evaluation is currently in progress; otherwise returns **disabled**.