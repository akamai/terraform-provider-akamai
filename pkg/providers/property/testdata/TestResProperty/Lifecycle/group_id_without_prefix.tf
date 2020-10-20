provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name        = "test property"
  product_id  = "prd_0"
  contract_id = "ctr_0"

  group_id = "0"
}
