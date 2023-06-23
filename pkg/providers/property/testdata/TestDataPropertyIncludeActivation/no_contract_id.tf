provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}


data "akamai_property_include_activation" "test" {
  group_id   = "group_321"
  include_id = "inc_1"
  network    = "STAGING"
}
