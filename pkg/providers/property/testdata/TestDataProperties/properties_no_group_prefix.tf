provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_properties" "akaproperties" {
  group_id    = "test"
  contract_id = "ctr_test"
}