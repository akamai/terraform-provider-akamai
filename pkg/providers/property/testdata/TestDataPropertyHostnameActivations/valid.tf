provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_hostname_activations" "activation" {
  contract_id = "1"
  group_id    = "1"
  property_id = "1"
}
