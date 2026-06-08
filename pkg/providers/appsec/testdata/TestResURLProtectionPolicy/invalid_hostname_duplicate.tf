resource "akamai_appsec_url_protection_policy" "test" {
  config_id          = 43007
  name               = "duplicate-hostnames"
  max_rate_threshold = 100

  hostname_paths = [
    {
      hostname = "example.com"
      paths    = ["/a"]
    },
    {
      hostname = "example.com"
      paths    = ["/a"]
    }
  ]
}


