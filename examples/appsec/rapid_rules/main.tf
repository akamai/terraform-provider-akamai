# This example presents a sample workflow for rapid rules.
# The `akamai_appsec_rapid_rules` resource enables and configures rapid rules.
# The `akamai_appsec_rapid_rules` data source returns information about rapid rules, including a name, action, action lock, attack group,
# exceptions, and the default action for new rapid rules and rapid ruleset status.
#
# To run this example:
#
# 1. Specify the path to your `.edgerc` file and the section header for the set of credentials to use.
#
# The defaults here expect the `.edgerc` at your home directory and use the credentials under the heading of `default`.
#
# 2. Make changes to the attribute values and in the `rule_definitions.json` file according to your needs.
#
# 3. Open a Terminal or shell instance and initialize the provider with `terraform init`. Then, run `terraform plan` to preview the changes and `terraform apply` to apply your changes.
#
# A successful operation enables and configures rapid rules and returns information about rapid rules, including a name, action, action lock, attack group,
# exceptions, and the default action and rapid ruleset status.

terraform {
  required_version = ">= 1.0"
  required_providers {
    akamai = {
      source  = "akamai/akamai"
      version = ">= 3.1.0"
    }
  }
}

provider "akamai" {
  edgerc         = "~/.edgerc"
  config_section = "default"
}

resource "akamai_appsec_rapid_rules" "rapid_rules" {
  config_id          = 273588
  security_policy_id = "4444_573430"
  default_action     = "akamai_managed"
  rule_definitions   = file("rule_definitions.json")
}

data "akamai_appsec_rapid_rules" "my_rapid_rules" {
  config_id          = akamai_appsec_rapid_rules.rapid_rules.config_id
  security_policy_id = akamai_appsec_rapid_rules.rapid_rules.security_policy_id
}

output "output_enabled" {
  value = data.akamai_appsec_rapid_rules.my_rapid_rules.enabled
}

output "output_default_action" {
  value = data.akamai_appsec_rapid_rules.my_rapid_rules.default_action
}

output "output_rapid_rules" {
  value = data.akamai_appsec_rapid_rules.my_rapid_rules.rapid_rules
}

output "output_text" {
  value = data.akamai_appsec_rapid_rules.my_rapid_rules.output_text
}

