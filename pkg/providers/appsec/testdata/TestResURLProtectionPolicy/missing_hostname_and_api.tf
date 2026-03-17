provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_url_protection_policy" "test" {
  config_id          = 43007
  name               = "URL Protection"
  max_rate_threshold = 195
}

