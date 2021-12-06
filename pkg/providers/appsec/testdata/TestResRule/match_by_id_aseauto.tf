provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_rule" "test" {
  config_id           = 43253
  security_policy_id  = "AAAA_81230"
  rule_id             = 12345
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


