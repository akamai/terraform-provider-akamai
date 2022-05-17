provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_rules" "rules" {
  contract_id = "ctr_2"
  property_id = "prp_2"
}
