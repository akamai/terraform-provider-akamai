provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_appsec_configuration_clone" "test" {
  create_from_config_id = 43253
  create_from_version   = 7
  name                  = "Test Configuratin"
  description           = "New configuration test"
  contract_id           = "C-1FRYVV3"
  group_id              = "64867"
  host_names            = ["rinaldi.sandbox.akamaideveloper.com", "sujala.sandbox.akamaideveloper.com"]
}

