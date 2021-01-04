provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_appsec_advanced_settings" "logging" {
  config_id = 43253
    version = 7
  logging  = <<-EOF
{
    "allowSampling": true,
    "cookies": {
        "type": "all"
    },
    "customHeaders": {
        "type": "exclude",
        "values": [
            "csdasdad"
        ]
    },
    "standardHeaders": {
        "type": "only",
        "values": [
            "Accept"
        ]
    }
}
EOF
}
