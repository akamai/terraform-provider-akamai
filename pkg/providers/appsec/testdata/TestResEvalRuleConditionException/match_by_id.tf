provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_eval_rule_condition_exception" "test" {
    config_id = 43253
    security_policy_id = "AAAA_81230"
    rule_id = 12345
   condition_exception  = <<-EOF
   {
    "conditions": [],
    "exception": {
        "headerCookieOrParamValues": [
            "abc"
        ],
        "specificHeaderCookieOrParamPrefix": {
            "prefix": "a*",
            "selector": "REQUEST_COOKIES"
        }
    }
}
EOF
}


