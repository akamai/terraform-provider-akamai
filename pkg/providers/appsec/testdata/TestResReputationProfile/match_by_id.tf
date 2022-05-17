provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_reputation_profile" "test" {
  config_id          = 43253
  reputation_profile = <<-EOF
{
    "name": "Web Attack Rep Profile",
    "description": "Reputation profile description",
    "context": "WEBATCK",
    "threshold": 5,
    "sharedIpHandling": "NON_SHARED",
    "condition": {
        "positiveMatch": true,
        "atomicConditions": [
            {
                "positiveMatch": true,
                "className": "AsNumberCondition",
                "value": [
                    "1"
                ]
            },
            {
                "positiveMatch": true,
                "nameWildcard": true,
                "valueWildcard": true,
                "className": "RequestCookieCondition",
                "nameCase": true,
                "name": "x-header"
            },
            {
                "positiveMatch": true,
                "valueWildcard": true,
                "className": "HostCondition",
                "host": [
                    "*.com"
                ]
            }
        ]
    }
}
EOF
}

