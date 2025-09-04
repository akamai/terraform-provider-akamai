provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_advanced_settings_ase_penalty_box" "test" {
  config_id      = 43253
  block_duration = 5
  qualification_exclusions {
    attack_groups = ["IN", "XSS"]
    rules         = [950002]
  }
}