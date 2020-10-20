provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property_rules" "rules" {
  contract_id = "1"
  group_id = "1"
  property_id = "1"
  rules = <<-EOF
{
        "name": "updated",
        "behaviors": [
            {
                "name": "beh_2"
            }
        ],
        "options": {
            "is_secure": true
        },
        "criteriaMustSatisfy": "all"
}
EOF
}
