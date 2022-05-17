provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_cloudlets_policy" "policy" {
  name          = "test_policy"
  cloudlet_code = "ER"
  description   = "test policy description"
  group_id      = "grp_123"
  match_rules   = <<-EOF
  [
  {
    "name": "r1",
    "type": "erMatchRule",
    "useRelativeUrl": "copy_scheme_hostname",
    "statusCode": 301,
    "redirectURL": "/ddd",
    "matchURL": "abc.com",
    "useIncomingQueryString": false,
    "useIncomingSchemeAndHost": true
  },
  {
    "name": "r3",
    "type": "erMatchRule",
    "matches": [
      {
        "matchType": "hostname",
        "matchValue": "3333.dom",
        "matchOperator": "equals",
        "caseSensitive": true,
        "negate": false
      }
    ],
    "useRelativeUrl": "copy_scheme_hostname",
    "statusCode": 307,
    "redirectURL": "/abc/sss",
    "useIncomingQueryString": false,
    "useIncomingSchemeAndHost": true
  }
]
EOF
}