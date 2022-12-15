provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_property_include" "test" {
  contract_id = "ctr_123"
  group_id    = "grp_123"
  name        = "test include"
  type        = "MICROSERVICES"
  rule_format = "v2022-06-28"
  rules       = "{}"
}