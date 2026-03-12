provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_url_protection_rule" "test" {
  config_id          = 43007
  name               = "URL Protection"
  max_rate_threshold = 100

  hostname_paths = [{
    hostname = "custom.com"
    paths    = ["/asd"]
  }]

  intelligent_load_shedding = {
    hits_per_sec = 95
    categories   = ["BOTS"]
  }
}
