provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_malware_policy" "test" {
  config_id      = 43253
  malware_policy = <<-EOF
{
  "name": "Updated FMS Configuration22",
  "description": "Malware scan configuration details",
  "hostnames": [
    "abc.com",
    "def.com",
    "xyz.com",
    "example.com"
  ],
  "paths": [
    "/base-path",
    "/test"
  ],
  "contentTypes": [
    {
      "name": "application/json",
      "encodedContentAttributes": [
        {
          "path": "image.imagePath1",
          "encoding": [
            "base64"
          ]
        }
      ]
    }
  ],
  "logFilename": false
}
EOF
}


