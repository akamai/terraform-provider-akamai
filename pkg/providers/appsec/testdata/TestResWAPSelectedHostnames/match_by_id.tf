provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_wap_selected_hostnames" "test" {
  config_id          = 43253
  security_policy_id = "AAAA_81230"
  protected_hosts    = ["rinaldi.sandbox.akamaideveloper.com"]
  evaluated_hosts    = ["sujala.sandbox.akamaideveloper.com"]
}

