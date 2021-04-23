provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_slow_post" "test" {
    config_id = 43253
    security_policy_id = "AAAA_81230"
    slow_rate_action = "alert"                        
    slow_rate_threshold_rate = 10
    slow_rate_threshold_period = 30
    duration_threshold_timeout = 20
}
