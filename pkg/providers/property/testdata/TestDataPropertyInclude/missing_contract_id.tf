provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}


data "akamai_property_include" "include" {
  group_id   = "grp_1"
  include_id = "inc_1"
}
