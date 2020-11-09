provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_aag_rule" "test" {
    config_id = 43253
    version = 7
    policy_id = "AAAA_81230"
    group_id  = "SQL"
   rules = <<-EOF
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


