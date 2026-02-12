provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_url_protection_rules_actions" "test" {
  config_id = 43253
}
