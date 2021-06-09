provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_rule_condition_exception" "test" {
    config_id = 43253
    security_policy_id = "AAAA_81230"
    rule_id = 12345
   condition_exception  = <<-EOF
   {
    "conditions": [
        {
            "type": "extensionMatch",
            "extensions": [
                "test"
            ],
            "positiveMatch": true
        },
        {
            "type": "filenameMatch",
            "filenames": [
                "test2"
            ],
            "positiveMatch": true
        },
        {
            "type": "hostMatch",
            "hosts": [
                "www.test.com"
            ],
            "positiveMatch": true
        },
        {
            "type": "ipMatch",
            "ips": [
                "123.123.123.123"
            ],
            "positiveMatch": true,
            "useHeaders": true
        },
        {
            "type": "uriQueryMatch",
            "caseSensitive": true,
            "name": "test3",
            "nameCase": false,
            "positiveMatch": true,
            "value": "test4",
            "wildcard": true
        },
        {
            "type": "requestHeaderMatch",
            "header": "referer",
            "positiveMatch": true,
            "value": "test5",
            "valueCase": false,
            "valueWildcard": false
        },
        {
            "type": "requestMethodMatch",
            "methods": [
                "GET"
            ],
            "positiveMatch": true
        },
        {
            "type": "pathMatch",
            "paths": [
                "/test6"
            ],
            "positiveMatch": true
        }
    ],
    "exception": {
        "headerCookieOrParamValues": [
            "test"
        ],
        "specificHeaderCookieOrParamNameValue": {
            "name": "test",
            "selector": "REQUEST_HEADERS",
            "value": "test"
        },
        "specificHeaderCookieOrParamNames": [
            {
                "names": [
                    "test"
                ],
                "selector": "REQUEST_HEADERS"
            },
            {
                "names": [
                    "test"
                ],
                "selector": "REQUEST_COOKIES"
            },
            {
                "names": [
                    "test"
                ],
                "selector": "ARGS"
            },
            {
                "names": [
                    "test"
                ],
                "selector": "JSON_PAIRS"
            },
            {
                "names": [
                    "test"
                ],
                "selector": "XML_PAIRS"
            }
        ],
        "specificHeaderCookieOrParamPrefix": {
            "prefix": "test",
            "selector": "REQUEST_HEADERS"
        }
    }
}
EOF
}

 
