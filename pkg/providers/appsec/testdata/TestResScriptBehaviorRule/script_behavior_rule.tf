provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_script_behavior_rule" "test" {
  config_id          = 40576
  security_policy_id = "ab_86352"

  rule_name        = "Block PCI Data Access"
  rule_description = "Prevents third-party scripts from accessing payment card data"
  activity_type    = "DATA_ACCESS"
  action           = "deny"

  page_condition = "MATCHES"
  page_values    = ["/checkout", "/billing/*"]

  script_type      = "THIRD_PARTY"
  script_condition = "MATCHES"
  script_values    = ["analytics.example.com/*"]

  activity_values_condition = "MATCHES"
  activity_values           = ["CREDIT_CARD", "PII"]
}