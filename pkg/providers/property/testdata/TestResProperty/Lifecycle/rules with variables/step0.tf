provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property" "test" {
  contract_id = "ctr_0"
  group_id    = "grp_0"
  name        = "test_property"
  product_id  = "prd_0"
  rules = jsonencode(
    {
      "rules" : {
        "behaviors" : [
          {
            "name" : "origin",
            "options" : {
              "cacheKeyHostname" : "REQUEST_HOST_HEADER",
              "compress" : true,
              "enableTrueClientIp" : true,
              "forwardHostHeader" : "REQUEST_HOST_HEADER",
              "hostname" : "test.domain",
              "httpPort" : 80,
              "httpsPort" : 443,
              "originCertificate" : "",
              "originSni" : true,
              "originType" : "CUSTOMER",
              "ports" : "",
              "trueClientIpClientSetting" : false,
              "trueClientIpHeader" : "True-Client-IP",
              "verificationMode" : "PLATFORM_SETTINGS"
            }
          }
        ],
        "children" : [
          {
            "behaviors" : [
              {
                "name" : "baseDirectory",
                "options" : {
                  "value" : "/smth/"
                }
              }
            ],
            "criteria" : [
              {
                "name" : "requestHeader",
                "options" : {
                  "headerName" : "Accept-Encoding",
                  "matchCaseSensitiveValue" : true,
                  "matchOperator" : "IS_ONE_OF",
                  "matchWildcardName" : false,
                  "matchWildcardValue" : false
                }
              }
            ],
            "name" : "change fwd path",
            "options" : {

            },
            "criteriaMustSatisfy" : "all"
          },
          {
            "behaviors" : [
              {
                "name" : "caching",
                "options" : {
                  "behavior" : "MAX_AGE",
                  "mustRevalidate" : false,
                  "ttl" : "1m"
                }
              }
            ],
            "name" : "caching",
            "options" : {

            },
            "criteriaMustSatisfy" : "any"
          }
        ],
        "comments" : "The behaviors in the Default Rule apply to all requests for the property hostname(s) unless another rule overrides the Default Rule settings.",
        "name" : "default",
        "options" : {

        },
        "variables" : [
          {
            "description" : "",
            "hidden" : true,
            "name" : "TEST_EMPTY_FIELDS",
            "sensitive" : false,
            "value" : ""
          },
          {
            "value" : "",
            "description" : null,
            "hidden" : true,
            "name" : "TEST_NIL_FIELD",
            "sensitive" : false
          }
        ]
      }
    }
  )

  hostnames {
    cert_provisioning_type = "DEFAULT"
    cname_from             = "from.test.domain"
    cname_to               = "to.test.domain"
    cname_type             = "EDGE_HOSTNAME"
  }
}
