provider "akamai" {
  edgerc = "~/.edgerc"
}

resource "akamai_property" "test" {
  name = "test property"
  group_id    = "grp_0"
  contract_id = "ctr_0"
  product_id  = "prd_0"
  # Fetch the newly created property
  depends_on = [
    akamai_property.test
  ]
  rules = "{\"rules\":{\"name\":\"update rule tree\"}}"

}
