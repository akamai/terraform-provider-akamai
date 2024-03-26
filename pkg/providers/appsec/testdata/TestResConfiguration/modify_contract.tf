provider "akamai" {
  edgerc        = "../../common/testutils/edgerc"
  cache_enabled = false
}

resource "akamai_appsec_configuration" "test" {
  name        = "Akamai Tools"
  description = "Akamai Tools"
  contract_id = "C-1FRYVV3_MODIFIED"
  group_id    = 64867
  host_names  = ["rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"]
}

