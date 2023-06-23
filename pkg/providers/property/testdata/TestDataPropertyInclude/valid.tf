provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}


data "akamai_property_include" "include" {
  contract_id = "ctr_1"
  group_id    = "grp_1"
  include_id  = "inc_1"
}
