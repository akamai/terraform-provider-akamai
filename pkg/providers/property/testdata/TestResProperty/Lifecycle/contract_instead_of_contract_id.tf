provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name       = "test property"
  contract   = "ctr_0"
  group_id   = "grp_0"
  product_id = "prd_0"
}
