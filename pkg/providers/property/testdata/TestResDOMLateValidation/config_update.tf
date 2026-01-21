provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_property_domainownership_late_validation" "test" {
  contract_id       = "ctr_1"
  group_id          = "grp_1"
  property_id       = "prp_123"
  version           = 2
  validation_method = "DNS_CNAME"
}