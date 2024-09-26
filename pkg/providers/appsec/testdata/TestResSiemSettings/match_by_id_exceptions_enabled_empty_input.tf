provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_siem_settings" "test" {
  config_id               = 43253
  enable_siem             = true
  enable_for_all_policies = false
  siem_id                 = 1
  security_policy_ids     = ["12345"]
  exceptions {
    rate           = ["alert"]
    bot_management = []
  }
}

