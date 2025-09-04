provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_associations" "test" {
  name = "abc"
}