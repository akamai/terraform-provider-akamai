provider "akamai" {
  edgerc = "../../test/edgerc"
}


data "akamai_include_activation" "test" {
  contract_id = "contract_123"
  group_id    = "group_321"
  include_id  = "inc_1"
}
