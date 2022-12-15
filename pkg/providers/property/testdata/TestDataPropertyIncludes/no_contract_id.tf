provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_includes" "test" {
  group_id = "group_321"
}