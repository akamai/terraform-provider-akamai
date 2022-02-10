provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

data "akamai_appsec_wap_selected_hostnames" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
}

