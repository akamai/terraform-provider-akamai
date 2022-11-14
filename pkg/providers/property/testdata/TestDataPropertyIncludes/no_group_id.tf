provider "akamai" {
  edgerc = "../../test/edgerc"
}

data "akamai_property_includes" "test" {
  contract_id = "contract_123"
}