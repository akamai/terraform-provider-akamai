provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_activations" "test" {
  ca_set_id = "12345"
  type      = "ACTIVATE"
}