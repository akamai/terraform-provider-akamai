provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_includes" "test" {
  group_id = "group_321"
}