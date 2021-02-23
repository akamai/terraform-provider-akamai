provider "akamai" {
  edgerc = "~/.edgerc"
}


resource "akamai_property" "prop" {
  name = "property_name"
  group_id = "grp_0"
  contract_id = "ctr_0"
  product_id = "prd_0"
}