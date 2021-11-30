provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_appsec_selected_hostnames" "test" {
  config_id = 43253
  hostnames = ["rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"]
  mode      = "REPLACE"
}

