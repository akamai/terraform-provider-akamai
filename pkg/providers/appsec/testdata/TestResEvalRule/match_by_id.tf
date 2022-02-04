provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

resource "akamai_appsec_eval_rule" "test" {
  config_id           = 43253
  security_policy_id  = "AAAA_81230"
  rule_id             = 12345
  rule_action         = "alert"
  condition_exception = <<-EOF
   {
        "conditions": [
            {
                "type": "extensionMatch",
                "extensions": [
                    "test"
                ],
                "positiveMatch": true
            }
        ],
        "exception": {
            "headerCookieOrParamValues": [
                "test"
            ]
        }
    }
EOF
}

