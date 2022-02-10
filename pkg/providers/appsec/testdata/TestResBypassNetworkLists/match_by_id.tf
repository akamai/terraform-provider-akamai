provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

resource "akamai_appsec_bypass_network_lists" "test" {
  config_id           = 43253
  bypass_network_list = ["888518_ACDDCKERS", "1304427_AAXXBBLIST"]
}

