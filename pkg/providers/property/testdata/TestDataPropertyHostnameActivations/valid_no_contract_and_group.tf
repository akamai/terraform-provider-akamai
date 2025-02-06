provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_hostname_activations" "activation" {
  property_id = "prp_1"
}
