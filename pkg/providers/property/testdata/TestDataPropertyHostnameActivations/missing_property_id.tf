provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_property_hostname_activations" "activation" {
  contract_id = "ctr_1"
  group_id    = "grp_1"
}
