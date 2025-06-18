provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_versions" "test" {
  id                   = "12345"
  include_certificates = false
}