provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_botman_content_protection_javascript_injection_rule" "test" {
  security_policy_id                              = "AAAA_81230"
  content_protection_javascript_injection_rule_id = "fake3f89-e179-4892-89cf-d5e623ba9dc7"
  content_protection_javascript_injection_rule = jsonencode(
    {
      "testKey" : "testValue3"
    }
  )
}
