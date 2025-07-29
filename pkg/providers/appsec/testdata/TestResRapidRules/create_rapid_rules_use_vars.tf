provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

# Variables
variable "default_action" {
  type    = string
  default = "deny"
}

variable "config_id" {
  type    = number
  default = 111111
}

variable "security_policy_id" {
  type    = string
  default = "2222_333333"
}

variable "rule_definitions_file" {
  type    = string
  default = "testdata/TestResRapidRules/RuleDefinitions.json"
}

resource "akamai_appsec_rapid_rules" "test" {
  config_id          = var.config_id
  security_policy_id = var.security_policy_id
  default_action     = var.default_action
  rule_definitions   = file(var.rule_definitions_file)
}
