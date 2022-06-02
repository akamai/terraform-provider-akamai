provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_properties" "akaproperties" {
  group_id    = "grp_test"
  contract_id = "ctr_test"
}