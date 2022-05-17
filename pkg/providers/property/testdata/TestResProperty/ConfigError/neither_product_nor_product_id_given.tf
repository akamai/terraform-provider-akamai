provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property" "test" {
  name        = "test_property"
  contract_id = "ctr_0"
  group_id    = "grp_0"
}
