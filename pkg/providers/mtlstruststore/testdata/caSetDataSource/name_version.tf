provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set" "test" {
  name    = "example-ca-set"
  version = 1
}