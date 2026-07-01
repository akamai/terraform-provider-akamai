provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

# intelligent_load_shedding is intentionally omitted (null) here.
# ValidateConfig should skip ILS validation when it is not set.
resource "akamai_appsec_url_protection_policy" "test" {
  config_id          = 43007
  name               = "URL Protection"
  max_rate_threshold = 195

  hostname_paths = [{
    hostname = "custom.com"
    paths    = ["/asd"]
  }]
}

