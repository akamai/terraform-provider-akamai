provider "akamai" {
  edgerc = "../../test/edgerc"
}


data "akamai_property_include_parents" "parents" {
  contract_id = "ctr_1"
  include_id  = "inc_1"
}
