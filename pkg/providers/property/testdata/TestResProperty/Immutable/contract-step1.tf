provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name       = "test property"
  group_id   = "grp_0"
  product_id = "prd_0"

  contract = "ctr_1"
}
