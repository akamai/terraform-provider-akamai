provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set" "test" {
  id   = "12345"
  name = "example-ca-set"
}