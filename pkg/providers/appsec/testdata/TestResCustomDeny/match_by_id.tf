provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_custom_deny" "test" {
  config_id   = 43253
  custom_deny = <<-EOF
{
    "name": "new_custom_deny",
    "description": "testing",
    "isPageUrl" : false,
    "parameters": [
        {
            "name": "response_status_code",
            "value": "403"
        },
        {
            "name": "prevent_browser_cache",
            "value": "true"
        },
        {
            "name": "response_content_type",
            "value": "application/json"
        },
        {
            "name": "response_body_content",
            "value": "new testing"
        }
    ]
}
EOF
}

