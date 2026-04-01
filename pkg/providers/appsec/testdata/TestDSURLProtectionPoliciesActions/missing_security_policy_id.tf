provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_appsec_url_protection_policies_actions" "test" {
  config_id = 43253
}
