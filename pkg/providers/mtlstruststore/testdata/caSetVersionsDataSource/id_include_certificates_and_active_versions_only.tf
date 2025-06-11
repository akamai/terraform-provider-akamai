provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

data "akamai_mtlstruststore_ca_set_versions" "test" {
  id                   = "12345"
  active_versions_only = true
  include_certificates = false
}