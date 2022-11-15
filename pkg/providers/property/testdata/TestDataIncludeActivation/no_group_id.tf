provider "akamai" {
  edgerc = "../../test/edgerc"
}


data "akamai_include_activation" "test" {
  contract_id = "contract_123"
  include_id  = "inc_321"
  network     = "STAGING"
}
