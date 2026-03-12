provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_waf_ruleset" "test" {
  config_id          = 111111
  security_policy_id = "2222_333333"

  attack_groups = [
    {
      attack_group        = "CMD"
      attack_group_action = "alert"
    },
    {
      attack_group        = "SQL"
      attack_group_action = "deny"
    },
    {
      attack_group        = "XSS"
      attack_group_action = "none"
    }
  ]
}
