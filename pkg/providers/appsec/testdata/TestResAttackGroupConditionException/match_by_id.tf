provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

resource "akamai_appsec_attack_group_condition_exception" "test" {
  config_id           = 43253
  security_policy_id  = "AAAA_81230"
  attack_group        = "SQL"
  condition_exception = <<-EOF
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

