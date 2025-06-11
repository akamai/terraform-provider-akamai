provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_versions" "test" {
  name = "test-ca-set-name"
}