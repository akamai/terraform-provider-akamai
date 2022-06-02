provider "akamai" {
  edgerc        = "../../test/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_rate_policy" "test" {
  config_id   = 43253
  rate_policy = <<-EOF
    {
    "matchType": "path",
    "type": "WAF",
    "name": "Test_Paths 3",
    "description": "AFW Test Extensions",
    "averageThreshold": 5,
    "burstThreshold": 10,
    "clientIdentifier": "ip",
    "useXForwardForHeaders": true,
    "requestType": "ClientRequest",
    "sameActionOnIpv6": false,
    "path": {
        "positiveMatch": true,
        "values": [
            "/login/",
            "/path/",
            "sec/"
        ]
    },
    "pathMatchType": "Custom",
    "pathUriPositiveMatch": true,
    "fileExtensions": {
        "positiveMatch": false,
        "values": [
            "3g2",
            "3gp",
            "aif",
            "aiff",
            "au",
            "avi",
            "bin",
            "bmp",
            "cab",
            "pcx"
        ]
    },
    "hostnames": [
        "www.ludin.org"
    ],
    "additionalMatchOptions": [
        {
            "positiveMatch": true,
            "type": "IpAddressCondition",
            "values": [
                "198.129.76.39"
            ]
        },
        {
            "positiveMatch": true,
            "type": "RequestMethodCondition",
            "values": [
                "GET"
            ]
        }
    ],
    "queryParameters": [
        {
            "name": "productId",
            "values": [
                "BUB_12",
                "SUSH_11"
            ],
            "positiveMatch": true,
            "valueInRange": false
        }
    ]
}
EOF
}


