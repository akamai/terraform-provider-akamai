provider "akamai" {
  edgerc = "../../common/testutils/edgerc"
}

resource "akamai_mtlstruststore_ca_set_activation" "test" {
  ca_set_id = "123456"
  version   = 2
  network   = "STAGING"
}
