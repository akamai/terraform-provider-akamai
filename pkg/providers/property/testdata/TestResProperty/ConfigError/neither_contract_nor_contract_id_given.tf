provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property" "test" {
  name       = "test_property"
  group_id   = "grp_0"
  product_id = "prd_0"
}
