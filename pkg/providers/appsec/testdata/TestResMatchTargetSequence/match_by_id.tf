provider "akamai" {
  edgerc = "~/.edgerc"
}



resource "akamai_appsec_match_target_sequence" "test" {
    config_id = 43253
    version = 7
    match_target_sequence = <<-EOF
    {
    "type": "website",
    "targetSequence": [
        {
            "targetId": 2052813,
            "sequence": 1
        },
        {
            "targetId": 2971336,
            "sequence": 2
        }
    ]
}
EOF

}
