provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_activation" "test" {
  id          = 321
  ca_set_name = "ab"
}