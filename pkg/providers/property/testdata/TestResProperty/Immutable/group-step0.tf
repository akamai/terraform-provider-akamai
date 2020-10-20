provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name        = "test property"
  contract_id = "ctr_0"
  product_id  = "prd_0"

  group = "grp_0"
}
