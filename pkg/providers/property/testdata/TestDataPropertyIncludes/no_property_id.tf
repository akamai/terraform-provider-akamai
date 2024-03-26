provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_includes" "test" {
  contract_id = "contract_123"
  group_id    = "group_321"
  parent_property {
    version = 2
  }
}