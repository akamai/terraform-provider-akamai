provider "akamai" {
  edgerc        = "~/.edgerc"
  cache_enabled = false
}

resource "akamai_appsec_configuration" "test" {
  name        = "Akamai Tools New"
  description = "Akamai Tools New"
  contract_id = "C-1FRYVV3"
  group_id    = 64867
  host_names  = ["rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"]
}

